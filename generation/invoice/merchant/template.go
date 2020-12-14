package merchant

import (
	template "./template"
)

func (self *Merchant) Template(title, sectionName, sectionDescription, content string) string {
	return `<html>
  <head>
    <title>` + title + `</title> 
    <style>` + template.Style() + `</style> 
  </head>
  <body>
    <header style="border-bottom: #ececec 24px solid;height:100px;margin-top:15px;padding-bottom:25px;">
      <div class="container">
        <div class="row">

          <div class="col-sm-1">
          </div>

          <div class="col-sm-6" style="padding-top:10px;">
            ` + HackwaveLogo().Size(100, 100).Style("float:left;margin-top:-5px;").HTML() + `
            <h1 style="font-size:32px;color:#2493eb;padding-bottom:0px;margin-bottom:2px;">` + self.Name + `</h1>
            <span style="margin-left:3px;padding-top:0px;margin-top:0px;color:#000;opacity:0.5;">` + self.Slogan + `</span>
          </div>

          <div class="col-sm-4" style="margin-top:20px;padding-bottom:10px;" >
            <div class="row">
              <div class="col-sm-6" style="text-align:right;vertical-align:middle;opacity:0.8;padding-bottom:5px;padding-top:5px;">
              <span><span style="font-size:25px;vertical-align:middle;">✉</span> Email</span>
            </div>
            <div class="col-sm-6" style="font-weight:100;text-align:left;vertical-align:middle;opacity:0.7;padding-bottom:5px;padding-top:11px;">
              <span>` + self.Contacts.Email() + `</span>
            </div>
          </div>

          <div class="row">
            <div class="col-sm-6" style="text-align:right;vertical-align:middle;opacity:0.8;padding-bottom:5px;padding-top:5px;">
              <span><span style="font-size:25px;vertical-align:middle;">✆</span> Phone</span>
            </div>
            <div class="col-sm-6" style="font-weight:100;text-align:left;vertical-align:middle;opacity:0.7;padding-bottom:5px;padding-top:11px;">
              <span>` + self.Contacts.Phone() + `</span>
            </div>
          </div>

          <div class="col-sm-1">
          </div>

        </div>

      </div>
    </div>
  </header> 


  <main>
    <div style="background-color:#ececec;margin-top:-10px;height:38px;font-color:#eee;font-size:17px;">
      <span style="float:left;margin-left:11%;">
        <span style="font-weight:100;">` + sectionName + `</span>
        <span style="margin-left:12px;">` + sectionDescription + `</span>
      </span>

      <span style="float:right;margin-right:11%;">
        <span style="margin-left:25px;margin-right:18px;"><strong>🌐</strong> https://` + self.WebsiteAddress + `</span>
        <strong>🔗</strong> https://github.com/` + self.Contacts.AccountAt("github.com") + `
      </span>
    </div>
	
    

    ` + content + `

      <div style="background-color:#ececec;height:132px;font-color:#eee;font-size:18px;margin-top:35px;text-align:left;padding-left:135px;padding-top:30px;padding-bottom:12px;padding-right:135px;line-height:45px;">
        <span>` + self.Name + ` developers contribute to a variety of projects including <strong>Go Langauge</strong>, <strong>Bitcoin</strong>, <strong>Debian Linux</strong>, <strong>Multiverse OS</strong>, <strong>Ruby</strong> and <strong>more</strong>...</span>
      </div>
    </main> 
  </body>
</html>`

}
