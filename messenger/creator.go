package messenger

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var headerTemplate = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD XHTML 1.0 Transitional //EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"> <html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office"> <head> <!--[if gte mso 9]> <xml> <o:OfficeDocumentSettings> <o:AllowPNG/> <o:PixelsPerInch>96</o:PixelsPerInch> </o:OfficeDocumentSettings> </xml> <![endif]--> <meta http-equiv="Content-Type" content="text/html; charset=UTF-8"> <meta name="viewport" content="width=device-width, initial-scale=1.0"> <meta name="x-apple-disable-message-reformatting"> <!--[if !mso]><!--> <meta http-equiv="X-UA-Compatible" content="IE=edge"><!--<![endif]--> <title></title> <style type="text/css"> @media only screen and (min-width: 620px) { .u-row { width: 600px !important; } .u-row .u-col { vertical-align: top; } .u-row .u-col-10p33 { width: 61.98px !important; } .u-row .u-col-33p33 { width: 199.98px !important; } .u-row .u-col-66p67 { width: 400.02px !important; } .u-row .u-col-89p67 { width: 538.02px !important; } .u-row .u-col-100 { width: 600px !important; } } @media (max-width: 620px) { .u-row .u-col h1 { text-align: center !important; line-height: 100% !important; } .u-row-container { max-width: 100% !important; padding-left: 0px !important; padding-right: 0px !important; } .u-row .u-col { min-width: 320px !important; max-width: 100% !important; display: block !important; } .u-row { width: 100% !important; } .u-col { width: 100% !important; } .u-col>div { margin: 0 auto; } } body { margin: 0; padding: 0; } table, tr, td { vertical-align: top; border-collapse: collapse; } p { margin: 0; } .ie-container table, .mso-container table { table-layout: fixed; } * { line-height: inherit; } a[x-apple-data-detectors='true'] { color: inherit !important; text-decoration: none !important; } table, td { color: #000000; } #u_body a { color: #0000ee; text-decoration: underline; } @media (max-width: 480px) { #u_content_heading_4 .v-text-align { text-align: center !important; } #u_column_8 .v-col-border { border-top: 0px solid transparent !important; border-left: 10px solid #ffffff !important; border-right: 10px solid #ffffff !important; border-bottom: 0px solid transparent !important; } #u_content_image_4 .v-container-padding-padding { padding: 40px 10px 0px !important; } #u_content_image_4 .v-src-width { width: 100% !important; } #u_content_image_4 .v-src-max-width { max-width: 100% !important; } #u_column_9 .v-col-border { border-top: 0px solid transparent !important; border-left: 10px solid #ffffff !important; border-right: 10px solid #ffffff !important; border-bottom: 1px solid #ffffff !important; } #u_content_text_10 .v-font-size { font-size: 13px !important; } #u_content_text_10 .v-text-align { text-align: center !important; } } </style> <!--[if !mso]><!--> <link rel="preconnect" href="https://fonts.googleapis.com"> <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin> <link href="https://fonts.googleapis.com/css2?family=Inter:wght@100..900&display=swap" rel="stylesheet"> <!--<![endif]--> </head> <body class="clean-body u_body" style="margin: 0;padding: 0;-webkit-text-size-adjust: 100%;background-color: #ffb258;color: #000000"> <!--[if IE]><div class="ie-container"><![endif]--> <!--[if mso]><div class="mso-container"><![endif]--> <table id="u_body" style="border-collapse: collapse;table-layout: fixed;border-spacing: 0;mso-table-lspace: 0pt;mso-table-rspace: 0pt;vertical-align: top;min-width: 320px;Margin: 0 auto;background-color: #ffb258;width:100%" cellpadding="0" cellspacing="0"> <tbody> <tr style="vertical-align: top"> <td style="word-break: break-word;border-collapse: collapse !important;vertical-align: top"> <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td align="center" style="background-color: #ffb258;"><![endif]--> <div class="u-row-container" style="padding: 0px;background-color: transparent"> <div class="u-row" style="margin: 0 auto;min-width: 320px;max-width: 600px;overflow-wrap: break-word;word-wrap: break-word;word-break: break-word;background-color: transparent;"> <div style="border-collapse: collapse;display: table;width: 100%;height: 100%;background-color: transparent;"> <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td style="padding: 0px;background-color: transparent;" align="center"><table cellpadding="0" cellspacing="0" border="0" style="width:600px;"><tr style="background-color: transparent;"><![endif]--> <!--[if (mso)|(IE)]><td align="center" width="61" class="v-col-border" style="width: 61px;padding: 0px;border-top: 0px solid transparent;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;border-radius: 0px;-webkit-border-radius: 0px; -moz-border-radius: 0px;" valign="top"><![endif]--> <div class="u-col u-col-10p33" style="max-width: 320px;min-width: 61.98px;display: table-cell;vertical-align: top;"> <div style="height: 100%;width: 100% !important;border-radius: 0px;-webkit-border-radius: 0px; -moz-border-radius: 0px;"> <!--[if (!mso)&(!IE)]><!--> <div class="v-col-border" style="box-sizing: border-box; height: 100%; padding: 0px;border-top: 0px solid transparent;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;border-radius: 0px;-webkit-border-radius: 0px; -moz-border-radius: 0px;"> <!--<![endif]--> <table style="font-family:'Inter',sans-serif;" role="presentation" cellpadding="0" cellspacing="0" width="100%" border="0"> <tbody> <tr> <td class="v-container-padding-padding" style="overflow-wrap:break-word;word-break:break-word;padding:10px;font-family:'Inter',sans-serif;" align="left"> <table width="100%" cellpadding="0" cellspacing="0" border="0"> <tr> <td class="v-text-align" style="padding-right: 0px;padding-left: 0px;" align="center"> <img align="center" border="0" src="https://kalendarz.krzysztofromanowski.pl/icon.png" width="192" height="192" alt="" title="" style="outline: none;text-decoration: none;-ms-interpolation-mode: bicubic;clear: both;display: inline-block !important;border: none;height: auto;float: none;width: 100%;max-width: 41.98px;" width="41.98" class="v-src-width v-src-max-width" /> </td> </tr> </table> </td> </tr> </tbody> </table> <!--[if (!mso)&(!IE)]><!--> </div><!--<![endif]--> </div> </div> <!--[if (mso)|(IE)]></td><![endif]--> <!--[if (mso)|(IE)]><td align="center" width="538" class="v-col-border" style="width: 538px;padding: 0px;border-top: 0px solid transparent;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;border-radius: 0px;-webkit-border-radius: 0px; -moz-border-radius: 0px;" valign="top"><![endif]--> <div class="u-col u-col-89p67" style="max-width: 320px;min-width: 538.02px;display: table-cell;vertical-align: top;"> <div style="height: 100%;width: 100% !important;border-radius: 0px;-webkit-border-radius: 0px; -moz-border-radius: 0px;"> <!--[if (!mso)&(!IE)]><!--> <div class="v-col-border" style="box-sizing: border-box; height: 100%; padding: 0px;border-top: 0px solid transparent;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;border-radius: 0px;-webkit-border-radius: 0px; -moz-border-radius: 0px;"> <!--<![endif]--> <table style="font-family:'Inter',sans-serif;" role="presentation" cellpadding="0" cellspacing="0" width="100%" border="0"> <tbody> <tr> <td class="v-container-padding-padding" style="overflow-wrap:break-word;word-break:break-word;padding:10px;font-family:'Inter',sans-serif;" align="left"> <!--[if mso]><table width="100%"><tr><td><![endif]--> <h1 class="v-text-align v-font-size" style="margin: 0px; line-height: 140%; text-align: left; word-wrap: break-word; font-size: 30px; font-weight: 300; letter-spacing: -2.4px"> Imprezy Speedcuberskie w Polsce</h1> <!--[if mso]></td></tr></table><![endif]--> </td> </tr> </tbody> </table> <!--[if (!mso)&(!IE)]><!--> </div><!--<![endif]--> </div> </div> <!--[if (mso)|(IE)]></td><![endif]--> <!--[if (mso)|(IE)]></tr></table></td></tr></table><![endif]--> </div> </div> </div>`
var footerTemplate = `<div class="u-row-container" style="padding: 0px;background-color: transparent"> <div class="u-row" style="margin: 0 auto;min-width: 320px;max-width: 600px;overflow-wrap: break-word;word-wrap: break-word;word-break: break-word;background-color: transparent;"> <div style="border-collapse: collapse;display: table;width: 100%;height: 100%;background-color: transparent;"> <!--[if (mso)|(IE)]><table width="100%" cellpadding="0" cellspacing="0" border="0"><tr><td style="padding: 0px;background-color: transparent;" align="center"><table cellpadding="0" cellspacing="0" border="0" style="width:600px;"><tr style="background-color: transparent;"><![endif]--> <!--[if (mso)|(IE)]><td align="center" width="600" class="v-col-border" style="background-color: #ffb258;width: 600px;padding: 0px;border-top: 0px solid transparent;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;border-radius: 0px;-webkit-border-radius: 0px; -moz-border-radius: 0px;" valign="top"><![endif]--> <div class="u-col u-col-100" style="max-width: 320px;min-width: 600px;display: table-cell;vertical-align: top;"> <div style="background-color: #ffb258;height: 100%;width: 100% !important;border-radius: 0px;-webkit-border-radius: 0px; -moz-border-radius: 0px;"> <!--[if (!mso)&(!IE)]><!--> <div class="v-col-border" style="box-sizing: border-box; height: 100%; padding: 0px;border-top: 0px solid transparent;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;border-radius: 0px;-webkit-border-radius: 0px; -moz-border-radius: 0px;"> <!--<![endif]--> <table id="u_content_text_10" style="font-family:'Inter',sans-serif;" role="presentation" cellpadding="0" cellspacing="0" width="100%" border="0"> <tbody> <tr> <td class="v-container-padding-padding" style="overflow-wrap:break-word;word-break:break-word;padding:20px 0px;font-family:'Inter',sans-serif;" align="left"> <div class="v-text-align v-font-size" style="font-size: 14px; color: #000000; line-height: 140%; text-align: center; word-wrap: break-word;"> <p style="line-height: 140%;">Wszystkie imprezy można obejrzeć pod adresem <strong style="line-height: 19.6px;"><a rel="noopener" href="https://kalendarz.krzysztofromanowski.pl/" target="_blank">kalendarz.krzysztofromanowski.pl</a></strong></p> </div> </td> </tr> </tbody> </table> <!--[if (!mso)&(!IE)]><!--> </div><!--<![endif]--> </div> </div> <!--[if (mso)|(IE)]></td><![endif]--> <!--[if (mso)|(IE)]></tr></table></td></tr></table><![endif]--> </div> </div> </div> <!--[if (mso)|(IE)]></td></tr></table><![endif]--> </td> </tr> </tbody> </table> <!--[if mso]></div><![endif]--> <!--[if IE]></div><![endif]--> </body> </html>`
var sectionHeaderTemplate = `
<div class="u-row-container section-header" style="padding: 0px;background-color: transparent">
	<div class="u-row" style="margin: 0 auto;min-width: 320px;max-width: 600px;overflow-wrap: break-word;word-wrap: break-word;word-break: break-word;background-color: transparent;">
		<div style="border-collapse: collapse;display: table;width: 100%;height: 100%;background-color: transparent;">
			<!--[if (mso)|(IE)]>
			<table width="100%" cellpadding="0" cellspacing="0" border="0">
				<tr>
					<td style="padding: 0px;background-color: transparent;" align="center">
						<table cellpadding="0" cellspacing="0" border="0" style="width:600px;">
							<tr style="background-color: transparent;">
							<![endif]--> <!--[if (mso)|(IE)]>
							<td align="center" width="600" class="v-col-border" style="background-color: #f9da61;width: 600px;padding: 0px;border-top: 3px solid #ffffff;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;border-radius: 0px;-webkit-border-radius: 0px; -moz-border-radius: 0px;" valign="top">
								<![endif]--> <div class="u-col u-col-100" style="max-width: 320px;min-width: 600px;display: table-cell;vertical-align: top;">
								<div style="background-color: #f9da61;height: 100%;width: 100% !important;border-radius: 0px;-webkit-border-radius: 0px; -moz-border-radius: 0px;">
									<!--[if (!mso)&(!IE)]><!-->
									<div class="v-col-border" style="box-sizing: border-box; height: 100%; padding: 0px;border-top: 3px solid #ffffff;border-left: 0px solid transparent;border-right: 0px solid transparent;border-bottom: 0px solid transparent;border-radius: 0px;-webkit-border-radius: 0px; -moz-border-radius: 0px;">
									<!--<![endif]-->
										<table id="u_content_heading_4" style="font-family:'Inter',sans-serif;" role="presentation" cellpadding="0" cellspacing="0" width="100%" border="0">
											<tbody>
												<tr>
													<td class="v-container-padding-padding" style="overflow-wrap:break-word;word-break:break-word;padding:15px 10px 15px 15px;font-family:'Inter',sans-serif;" align="left">
													<!--[if mso]>
														<table width="100%">
															<tr>
																<td>
																<![endif]-->
																	<h1 class="v-text-align v-font-size" style="margin: 0px; color: #000000; line-height: 140%; text-align: left; word-wrap: break-word; font-size: 22px; font-weight: 600;">
	{{.}}
																	</h1>
																<!--[if mso]>
																</td>
															</tr>
														</table>
													<![endif]-->
													</td>
												</tr>
											</tbody>
										</table>
										<!--[if (!mso)&(!IE)]><!-->
										</div>
									<!--<![endif]-->
									</div>
								</div>
								<!--[if (mso)|(IE)]>
								</td>
							<![endif]-->
							<!--[if (mso)|(IE)]>
							</tr>
						</table>
					</td>
				</tr>
			</table>
			<![endif]-->
		</div>
	</div>
</div>`
var eventImageTemplate = `<img src="https://raw.githubusercontent.com/castus/castus.github.io/refs/heads/master/assets/events/{{.}}.png" class="event-image-{{.}}" width="20" height="20" alt="{{.}}" style="vertical-align: bottom;" />`
var itemTemplate = `
<div class="u-row-container competition-item" style="padding: 0px;background-color: transparent">
	<div class="u-row" style="margin: 0 auto;min-width: 320px;max-width: 600px;overflow-wrap: break-word;word-wrap: break-word;word-break: break-word;background-color: transparent;">
		<table border="1" cellpadding="10px" style="margin: 0;border-collapse: collapse;width: 100%;background: white;border: 10px solid #F9DA61;border-top: 1px solid #F9DA61;border-bottom: 5px solid #F9DA61;">
			<tbody>
				<tr>
					<td style="width: 100px;text-align: center;border-right: 1px solid white;">
						<img align="center" border="0" src="{{ .LogoURL }}" alt="logo" title="logo" style="outline: none;text-decoration: none;-ms-interpolation-mode: bicubic;clear: both;display: inline-block !important;border: none;height: auto;float: none;width: 100%;max-width: 140px;" width="140" class="v-src-width v-src-max-width" />
					</td>
					<td>
						<table id="u_content_heading_3" style="font-family:'Inter',sans-serif;" role="presentation" cellpadding="0" cellspacing="0" width="100%" border="0">
							<tbody>
								<tr>
									<td class="v-container-padding-padding" style="overflow-wrap:break-word;word-break:break-word;font-family:'Inter',sans-serif;" align="left">
										<!--[if mso]>
										<table width="100%">
											<tr>
												<td>
													<![endif]-->
													<h1 class="v-text-align v-font-size" style="letter-spacing: -0.8px; margin: 0px; color: #000000; line-height: 140%; text-align: left; word-wrap: break-word; font-family:'Inter',sans-serif; font-size: 22px; font-weight: 500;">
														{{ if .URL }}
															<a href="{{ .URL }}" class="header-with-link">
														{{ end }}
														{{ .Header }}
														{{ if .HasWCA }}
																<img src="https://www.speedcubing.pl/images/wca_small_logo.png" width="20" height="20" class="header-wca" />
														{{ end }}
														{{ if .URL }}
															</a>
														{{ end }}
													</h1>
													<!--[if mso]>
												</td>
											</tr>
										</table>
										<![endif]-->
									</td>
								</tr>
							</tbody>
						</table>
						<table id="u_content_text_6" style="font-family:'Inter',sans-serif;" role="presentation" cellpadding="0" cellspacing="0" width="100%" border="0">
							<tbody>
								<tr>
									<td class="v-container-padding-padding" style="overflow-wrap:break-word;word-break:break-word;/* padding:0px 50px 10px 10px; */font-family:'Inter',sans-serif;" align="left">
										<div class="v-text-align v-font-size" style="font-size: 14px; color: #000000; line-height: 140%; text-align: left; word-wrap: break-word;">
											<p style="margin: 0 0 3px">{{.Date}}</p>
											<p style="margin: 0 0 3px">{{.Place}}	{{if ne .Place "zawody online"}} <small style="color: #8d8d8d" class="distance-info">{{.Distance}}, {{.Duration}} jazdy autem</small>{{end}}</p>
											{{if ne (len .Events) 0}}
												<p style="margin: 0 0 3px" class="events">Konkurencje: {{$that := .}}{{range .Events}}{{call $that.Function .}} {{end}}</p>
											{{end}}

											{{if ne (len .MainEvent) 0}}
												<p style="margin: 0 0 3px" class="main-event">Konkurencja główna: {{call .Function .MainEvent}}</p>
											{{end}}

											{{if ne .CompetitorLimit 0}}
												<p style="margin: 0 0 3px" class="registered">Zarejestrowanych:
													{{if eq .Registered .CompetitorLimit}}
														<span style="color:#f14561" class="registered-limit">{{.Registered}}/{{.CompetitorLimit}}</span>
													{{else}}
														{{.Registered}}/{{.CompetitorLimit}}
													{{end}}
												</p>
											{{end}}

											<p style="margin: 0 0 3px"><a href="mailto:{{.ContactURL}}">{{.ContactName}} ({{.ContactURL}})</a></p>
										</div>
									</td>
								</tr>
							</tbody>
						</table>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>`

func PrepareMessage(added []MessengerDTO, passed []MessengerDTO, changed []MessengerDTO) string {
	message := []string{headerTemplate}
	if len(added) > 0 {
		message = append(message, prepareAdded(added))
	}
	if len(changed) > 0 {
		message = append(message, prepareChanged(changed))
	}
	if len(passed) > 0 {
		message = append(message, preparePassed(passed))
	}
	message = append(message, footerTemplate)

	m := strings.Join(message, "\n")
	m = strings.ReplaceAll(m, "\n", "")
	m = strings.ReplaceAll(m, "  ", "")

	return m
}

func prepareAdded(competitions []MessengerDTO) string {
	var message []string
	message = append(message, parseTemplate(sectionHeaderTemplate, "Imprezy dodane"))
	message = append(message, itemsHTML(competitions))

	return strings.Join(message, "\n")
}

func preparePassed(competitions []MessengerDTO) string {
	var message []string
	message = append(message, parseTemplate(sectionHeaderTemplate, "Imprezy minione"))
	message = append(message, itemsHTML(competitions))

	return strings.Join(message, "\n")
}

func prepareChanged(competitions []MessengerDTO) string {
	var message []string
	message = append(message, parseTemplate(sectionHeaderTemplate, "Imprezy zmienione"))
	message = append(message, itemsHTML(competitions))

	return strings.Join(message, "\n")
}

func itemsHTML(items []MessengerDTO) string {
	var message []string
	for _, item := range items {
		message = append(message, formattedItemAsHTML(item))
	}

	return strings.Join(message, "\n")
}

type MessengerDTO struct {
	LogoURL         string
	Name            string
	URL             string
	HasWCA          bool
	Date            string
	Distance        string
	Duration        string
	Place           string
	Events          []string
	MainEvent       string
	CompetitorLimit int
	Registered      int
	ContactURL      string
	ContactName     string
}

type itemData struct {
	LogoURL         string
	Name            string
	URL             string
	HasWCA          bool
	Date            string
	Distance        string
	Duration        string
	Place           string
	Events          []string
	MainEvent       string
	CompetitorLimit int
	Registered      int
	ContactURL      string
	ContactName     string
	Header          string
	Function        func(string) string
}

func formattedItemAsHTML(c MessengerDTO) string {
	data := itemData{
		LogoURL:         c.LogoURL,
		URL:             c.URL,
		Header:          makeHeader(c.Name),
		HasWCA:          c.HasWCA,
		Date:            c.Date,
		Place:           c.Place,
		Distance:        c.Distance,
		Duration:        c.Duration,
		Events:          c.Events,
		Function:        eventHTML,
		MainEvent:       c.MainEvent,
		CompetitorLimit: c.CompetitorLimit,
		Registered:      c.Registered,
		ContactURL:      c.ContactURL,
		ContactName:     c.ContactName,
	}
	return parseTemplate(itemTemplate, data)
}

func makeHeader(title string) string {
	header := fmt.Sprintf("%s", title)
	header = strings.ReplaceAll(header, "2024", "")
	header = strings.Trim(header, " ")

	return header
}

func eventHTML(event string) string {
	return parseTemplate(eventImageTemplate, event)
}

func parseTemplate(tmpl string, data any) string {
	var output bytes.Buffer
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	err = t.ExecuteTemplate(&output, "template", data)
	if err != nil {
		panic(err)
	}

	return output.String()
}
