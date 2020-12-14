package merchant

import (
	"fmt"

	currency "../currency"
	order "./order"
)

func (self *Merchant) Invoice(o *order.Order) string {
	return self.Template(
		self.Name+` Invoice `+o.InvoiceID(),
		"INVOICE",
		o.InvoiceID(),
		self.InvoiceInformation(o)+self.OrderInformation(o),
	)
}

func (self *Merchant) InvoiceInformation(o *order.Order) string {
	return `<div style="margin-top:35px;margin-bottom:35px;">
        <table style="margin-top:15px;width:80%;margin-left:10%;"> 
          <thead>
            <tr>
              <th style="font-weight:100; border-top:solid 1px #ccc;border-left:solid 1px #ccc;text-align:right;border-bottom:solid 1px #ccc;">
                INVOICE
              </th>
              <th style="font-weight:700; border-top:solid 1px #ccc;text-align:left;border-bottom:solid 1px #ccc;border-left:solid 1px #ccc;">
                ` + o.InvoiceID() + `
              </th>
              <th style="font-weight:100; border-top:solid 1px #ccc;text-align:right;border-bottom:solid 1px #ccc;">
                CREATED ON
              </th>
              <th style="font-weight:700; border-top:solid 1px #ccc;border-left:solid 1px #ccc;text-align:left;border-bottom:solid 1px #ccc;">
                ` + fmt.Sprintf("%v", o.CreatedAt()) + `
              </th>
            </tr>
        </thead>
      <tbody>
        <tr>
        <td style="text-align:right;width:25%;border-left:solid 1px #ccc;">
            <strong>Customer Name</strong>
          </td>
          <td style="text-align:left;width:25%;">
            ` + o.Customer.FullName + `
          </td>
          <td style="text-align:right;width:25%;">
            <strong>Company Name</strong>
          </td>
          <td style="text-align:left;width:25%;border-right:solid 1px #ccc;">
          ` + o.Customer.BusinessName() + `
          </td>
        </tr>
		
        <tr>
          <td style="text-align:right;width:25%;border-left:solid 1px #ccc;border-bottom:solid 1px #ccc;">
            <strong>LocalBitcoins.com Account</strong>
          </td>
          <td style="text-align:left;width:25%;border-bottom:solid 1px #ccc;">
			` + o.Data["localbitcoins.com"]["username"] + `
          </td>
          <td style="text-align:right;width:25%;border-bottom:solid 1px #ccc;">
            <strong>LocalBitcoins Reference</strong>
          </td>
          <td style="text-align:left;width:25%;border-right:solid 1px #ccc;border-bottom:solid 1px #ccc;">
			` + o.Data["localbitcoins.com"]["reference"] + `
          </td>
        </tr>

        <tr>
          <td style="text-align:right;width:25%;border-left:solid 1px #ccc;border-bottom:solid 1px #ccc;">
            <strong>Mailing Address</strong>
          </td>
          <td style="text-align:left;width:25%;border-bottom:solid 1px #ccc;">
            ` + o.Customer.Address.Format() + `
          </td>
          <td style="text-align:right;width:25%;border-bottom:solid 1px #ccc;">
            <strong>Phone</strong>
          </td>
          <td style="text-align:left;width:25%;border-right:solid 1px #ccc;border-bottom:solid 1px #ccc;">
            ` + o.Customer.Contacts.Phone() + `
          </td>
        </tr>


      </tbody>
     </table>
    </div>`
}

func (self *Merchant) OrderLineItem(lineItem *order.LineItem) string {
	return `<tr>
              <td style="text-align:center;border-left: solid 1px #ccc;">` + lineItem.ID + `</td>
              <td style="text-align:center;"><strong>` + lineItem.Name + `</strong></td>
              <td>` + lineItem.Description + `</td>
              <td style="text-align:center;">x` + fmt.Sprintf("%.4f", lineItem.Quantity) + `</td>
              <td style="text-align:center;border-right: solid 1px #ccc;">` + fmt.Sprintf("%.2f", lineItem.Price) + ` <strong>` + lineItem.Currency.String() + `</strong></td>
            </tr>`
}

func (self *Merchant) OrderLineItems(o *order.Order) (lineItems string) {
	for _, lineItem := range o.LineItems {
		lineItems += self.OrderLineItem(lineItem)
	}
	return lineItems
}

// TODO Add the conversion from UYU currencies
func (self *Merchant) OrderTotal(o *order.Order) string {
	if o.Currency == currency.UYU {
		fmt.Println("Order paid with UYU, converting to USD for clairity")
		return `<tr>
             <td colspan="3">
             </td>
             <td style="text-align:center;border-left: 1px solid #ccc;border-bottom: 1px solid #ccc;border-top: 1px solid #ccc;">
            <strong>UYU Total</strong>
             </td>
             <td style="text-align:center;border-right: 1px solid #ccc;border-bottom: 1px solid #ccc;border-top: 1px solid #ccc;">
               ` + fmt.Sprintf("%.2f", o.Total()) + ` <strong>` + o.Currency.String() + `</strong>
             </td>
           </tr>
		    <tr>
             <td colspan="3">
             </td>
             <td style="text-align:center;border-left: 1px solid #ccc;border-bottom: 1px solid #ccc;border-top: 1px solid #ccc;">
            <strong>Total</strong>
             </td>
             <td style="text-align:center;border-right: 1px solid #ccc;border-bottom: 1px solid #ccc;border-top: 1px solid #ccc;">
               ` + fmt.Sprintf("%.2f", ConvertUYUtoUSD(o.Total())) + ` <strong>` + self.Currency.String() + `</strong>
             </td>
           </tr>`
	} else {
		return `<tr>
             <td colspan="3">
             </td>
             <td style="text-align:center;border-left: 1px solid #ccc;border-bottom: 1px solid #ccc;border-top: 1px solid #ccc;">
            <strong>Total</strong>
             </td>
             <td style="text-align:center;border-right: 1px solid #ccc;border-bottom: 1px solid #ccc;border-top: 1px solid #ccc;">
               ` + fmt.Sprintf("%.2f", o.Total()) + ` <strong>` + o.Currency.String() + `</strong>
             </td>
           </tr>`
	}
}

func (self *Merchant) OrderInformation(o *order.Order) string {
	return `<table style="margin-top:15px;width:80%;margin-left:10%;"> 
    <thead>
      <th style="width:12%;text-align:center;border-left:solid 1px #ccc;border-top:solid 1px #ccc;border-bottom:solid 1px #ccc;">ID</th>
      <th style="width:15%;text-align:center;border-top:solid 1px #ccc;border-bottom:solid 1px #ccc;">Name</th>
      <th style="width:51%;border-top:solid 1px #ccc;border-bottom:solid 1px #ccc;">Description</th>
      <th style="width:10%;text-align:center;border-top:solid 1px #ccc;border-bottom:solid 1px #ccc;">Qty.</th>
      <th style="width:12%;text-align:center;border-top:solid 1px #ccc;border-right:solid 1px #ccc;border-bottom:solid 1px #ccc;">Price</th>
    </thead>
    <tbody>
      ` + self.OrderLineItems(o) + `
      ` + self.OrderTotal(o) + `
    </tbody>
  </table>`
}
