package merchant

import (
	"fmt"
)

func (self *Merchant) YearReport(year int) string {
	title := self.Name + fmt.Sprintf("%v", year) + ` Sales Report`
	return self.Template(
		title,
		"SALES REPORT",
		fmt.Sprintf("%v", year),
		self.SalesReport(),
	)

}

func (self *Merchant) SalesReport() string {
	return `<table style="margin-top:15px;width:80%;margin-left:10%;"> 
    <thead>
    <th style="width:12%;text-align:center;border-left:solid 1px #ccc;border-top:solid 1px #ccc;border-bottom:solid 1px #ccc;">Created At</th>
        <th style="width:15%;text-align:center;border-top:solid 1px #ccc;border-bottom:solid 1px #ccc;">ID</th>
        <th style="width:51%;border-top:solid 1px #ccc;border-bottom:solid 1px #ccc;">Customer</th>
        <th style="width:12%;text-align:center;border-top:solid 1px #ccc;border-right:solid 1px #ccc;border-bottom:solid 1px #ccc;">Total</th>
        </thead>
        <tbody>
		` + self.OrderTableContent() + `
		    <tr>
             <td colspan="2">
             </td>
             <td style="text-align:center;border-left: 1px solid #ccc;border-bottom: 1px solid #ccc;border-top: 1px solid #ccc;">
            <strong>Year Total</strong>
             </td>
			 <td style="text-align:center;border-right: 1px solid #ccc;border-bottom: 1px solid #ccc;border-top: 1px solid #ccc;">
			   <span style="color:#2ECC40;">` + fmt.Sprintf("%.2f", self.OrdersTotal()) + `</span> <strong>` + self.Currency.String() + `</span> 
             </td>
           </tr>
        </tbody>
    </table>`
}

//<td>` + order.Customer.Name() + `</td>
func (self *Merchant) OrderTableContent() (tableContent string) {
	var customerName string
	for _, order := range self.History.Orders {
		customerName = order.Customer.FullName
		if order.Customer.IsBusiness() {
			customerName += ` representative for ` + order.Customer.Business
		}

		tableContent += `<tr>
<td style="text-align:center;font-weight:100;">` + order.CreatedAt() + `</td>
			<td style="text-align:center;"><a style="font-style:none;font-weight:200;font-color:#000;" href="/invoice/by-index/` + fmt.Sprintf("%v", order.ID) + `">` + order.InvoiceID() + `</a></td>
			<td>` + customerName + ` </td>
            <td style="text-align:center;">` + fmt.Sprintf("%.2f", order.Total()) + ` <strong>` + self.Currency.String() + `</strong></td>
		</tr>`
	}
	return tableContent
}
