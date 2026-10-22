# Portfolio Management

Módulo escrito en Go para administrar un portafolio de inversiones y calcular las operaciones de compra o venta necesarias para alcanzar una distribución objetivo.

## Responsabilidades

1. **Stock:** representa una acción mediante su símbolo, cantidad y precio actual.
2. **Portfolio:** administra las posiciones del inversionista y calcula su valor total.
3. **Allocation:** define el porcentaje objetivo de cada acción dentro del portafolio.
4. **Rebalance:** genera una propuesta de operaciones para alcanzar la distribución objetivo.

## Ejemplo

Supongamos que el portafolio tiene la siguiente distribución:

| Acción | Valor actual | Distribución actual | Objetivo |
|---|---:|---:|---:|
| META | $7,000 | 70% | 40% |
| AAPL | $3,000 | 30% | 60% |
| **Total** | **$10,000** | **100%** | **100%** |

El rebalanceo propone:

- Vender $3,000 de META.
- Comprar $3,000 de AAPL.

El resultado sería un portafolio con $4,000 en META y $6,000 en AAPL, alcanzando la distribución objetivo 40/60.

## Decisiones de diseño

- Se utilizan `struct` siguiendo las convenciones de Go.
- El cálculo de rebalanceo está separado de la ejecución de operaciones bursátiles.
- Se trabaja con cantidades de acciones y precios, no solamente con porcentajes.
- `Rebalance` devuelve una propuesta de operaciones sin modificar el portafolio.
- Se permiten acciones fraccionarias.
- Se utiliza `float64` para este prototipo. En una aplicación financiera real convendría utilizar una representación decimal precisa para importes monetarios.
- Las pruebas unitarias usan table-driven tests, el patrón Arrange-Act-Assert y `testify/assert`.

## Estructura del proyecto

```text
portfolio-management/
├── cmd/
│   └── example/
│       └── main.go
├── portfolio/
│   ├── portfolio.go
│   ├── rebalance.go
│   ├── rebalance_test.go
│   └── stock.go
├── go.mod
├── go.sum
└── README.md
```

## Modelos

### Stock

```go
type Stock struct {
    Symbol   string
    Quantity float64
    Price    float64
}
```

`MarketValue` calcula el valor actual de una posición:

```text
MarketValue = Quantity × Price
```

### Portfolio

```go
type Portfolio struct {
    Stocks      []Stock
    Allocations map[string]float64
}
```

Las asignaciones utilizan valores decimales entre `0` y `1`:

```go
allocations := map[string]float64{
    "META": 0.40,
    "AAPL": 0.60,
}
```

### Operation

Cada operación propuesta contiene el símbolo, tipo de operación, cantidad de acciones, precio y valor total:

```go
type Operation struct {
    Symbol   string
    Type     OperationType
    Quantity float64
    Price    float64
    Value    float64
}
```

Los tipos posibles son `BUY` y `SELL`.

## Algoritmo de rebalanceo

Para cada acción se realizan los siguientes cálculos:

```text
TargetValue = TotalPortfolioValue × Allocation
Difference  = TargetValue - CurrentValue
Quantity    = |Difference| ÷ CurrentPrice
```

- Si `Difference` es positiva, se propone una compra.
- Si `Difference` es negativa, se propone una venta.
- Si la diferencia está dentro de la tolerancia definida, no se genera ninguna operación.

Las operaciones se ordenan alfabéticamente por símbolo para producir resultados deterministas.

## Validaciones

Antes de calcular el rebalanceo se verifica que:

- Exista al menos una asignación objetivo.
- Las asignaciones sumen `1` (100%).
- Cada asignación se encuentre entre `0` y `1`.
- Los símbolos no estén vacíos ni duplicados.
- Las cantidades no sean negativas.
- Los precios sean mayores que cero.
- Cada posición tenga una asignación objetivo.
- Cada asignación tenga una posición con precio disponible.

## Uso

```go
p := portfolio.Portfolio{
    Stocks: []portfolio.Stock{
        {Symbol: "META", Quantity: 70, Price: 100},
        {Symbol: "AAPL", Quantity: 30, Price: 100},
    },
    Allocations: map[string]float64{
        "META": 0.40,
        "AAPL": 0.60,
    },
}

operations, err := p.Rebalance()
if err != nil {
    log.Fatal(err)
}
```

Para ejecutar el ejemplo completo:

```bash
go run ./cmd/example
```

Salida esperada:

```text
BUY 30.00 shares of AAPL at $100.00 ($3000.00)
SELL 30.00 shares of META at $100.00 ($3000.00)
```

## Pruebas

Ejecuta todas las pruebas con:

```bash
go test ./...
```

Las pruebas cubren:

- Portafolios balanceados y desbalanceados.
- Portafolios vacíos.
- Configuraciones válidas e inválidas.
- Cálculo del valor total.
- Inmutabilidad del portafolio durante el rebalanceo.
