1. Análisis de requerimientos
El problema consiste en implementar un módulo que permita administrar un portafolio de inversiones y calcular qué acciones comprar o vender para alcanzar una distribución objetivo.
Identifico cuatro responsabilidades principales:
1. Stock: representar una acción, su símbolo y su precio actual.
2. Portfolio: administrar las acciones que posee el inversionista y sus cantidades.
3. Allocation: definir qué porcentaje del valor total debe representar cada acción.
4. Rebalance: calcular las operaciones necesarias para aproximarse a la distribución objetivo.
Ejemplo
Supongamos que tenemos este portafolio:
Acción	Valor actual	Distribución actual	Objetivo
META	$7,000	70%	40%
AAPL	$3,000	30%	60%
Total	$10,000	100%	100%
El resultado esperado sería:
- Vender $3,000 de META.
- Comprar $3,000 de AAPL.
Así tendríamos $4,000 en META y $6,000 en AAPL, alcanzando la distribución 40/60.
2. Decisiones de diseño
Para la primera versión propongo:
- Utilizar struct en lugar de clases, siguiendo las convenciones de Go.
- Separar el cálculo de rebalanceo de la ejecución de operaciones bursátiles.
- Trabajar con cantidades de acciones y precios, no únicamente porcentajes.
- Hacer que Rebalance devuelva una propuesta de operaciones sin modificar el portafolio.
- Implementar pruebas unitarias con testing, sin dependencias adicionales.
Hay una decisión importante: el ejercicio no especifica si se permiten acciones fraccionarias. Esto afecta directamente al cálculo. Para empezar podemos permitir fracciones y después implementar una versión que utilice exclusivamente acciones enteras.




3. Estructura inicial del proyecto
Propongo comenzar con esta estructura, sin complicarnos todavía con arquitectura hexagonal:
portfolio-management/
├── go.mod
├── portfolio/
│   ├── stock.go
│   ├── portfolio.go
│   ├── rebalance.go
│   └── rebalance_test.go
├── cmd/
│   └── example/
│       └── main.go
└── README.md


Primer modelo: Stock
package portfolio

type Stock struct {
    Symbol   string
    Quantity float64
    Price    float64
}

// CurrentPrice returns the last available market price.
func (s Stock) CurrentPrice() float64 {
    return s.Price
}

// MarketValue calculates the total value of a position.
func (s Stock) MarketValue() float64 {
    return s.Quantity * s.CurrentPrice()
}


Segundo modelo: Portfolio
package portfolio

type Portfolio struct {
    Stocks      []Stock
    Allocations map[string]float64
}

// TotalValue returns the current market value of the portfolio.
func (p Portfolio) TotalValue() float64 {
    var total float64

    for _, stock := range p.Stocks {
        total += stock.MarketValue()
    }

    return total
}


Aquí Allocations utiliza porcentajes decimales, por ejemplo:
allocations := map[string]float64{
    "META": 0.40,
    "AAPL": 0.60,
}


4. Algoritmo de rebalanceo
El cálculo principal es:
\[
\text{TargetValue} = \text{TotalPortfolioValue} \times \text{Allocation}
\]

\[
\text{Difference} = \text{TargetValue} - \text{CurrentValue}
\]

\[
\text{Quantity} = \frac{|\text{Difference}|}{\text{CurrentPrice}}
\]

Si la diferencia es positiva, debemos comprar. Si es negativa, debemos vender.
Antes de implementarlo, tendremos que contemplar validaciones como porcentajes que no sumen 100%, precios inválidos, acciones sin precio y cantidades negativas.
También debemos considerar que float64 sirve para un prototipo, pero en aplicaciones financieras reales conviene usar una representación decimal precisa para los importes monetarios.
5. Orden de implementación
Mi propuesta es desarrollar el ejercicio en tres etapas:
1. Modelado: definir Stock, Portfolio y las reglas de validación.
2. Rebalanceo: implementar el cálculo de compras y ventas sin ejecutar transacciones.
3. Testing: cubrir portafolios balanceados, desbalanceados, vacíos y configuraciones inválidas.