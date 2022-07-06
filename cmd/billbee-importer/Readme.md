# Billbee Importer

1.) Run the Application of your OS. (Either Linux, MacOs or Windows)

1.1) On Linux / MacOs start the Application in the Terminal.

2.) Provide settings in the `config.dist.yml` file.

3.) Provide stockx sales in the `stockx.csv` file. See `stockx-example.csv` to see the csv structure.

4.) Provide alias sales in the `alias.csv` file. See `alias-example.csv` to see the csv structure.

5.) Please test with a single row first and check the generated result before applying a csv with hundreds of rows. 

### Important
* You need to provide your email smtp settings. See `config.dist.yml`
* The stockx.csv file has to be named `stockx.csv`
* The alias.csv file has to be named `alias.csv`
* Please read the billbee docs!

### Setup Billbee for Stockx
* go to https://app.billbee.io/app_v2/settings/shops -> "new" -> "Manueller Shop" at the bottom of the list
* Eigenschaften: -> name your shop
* Umsatzsteuer: -> set Umsatzsteuer to "innergemeinschaftliche Leistung (Nettopreise)
* Layouts und Nummernkreise: set "Layout Rechnung" to "StockX"
* (to create StockX Layout, go to https://app.billbee.io/app_v2/settings -> Auftragsdokumente. Create a new Layout for your invoices. Go to "Fussbereich", put "Das Versanddatum entspricht dem Lieferdatum steuerfreie innergemeinschaftliche Lieferung" into Fußbereich Spalte 1.)
* That's it for the basics. If you want to add some spice, feel free to do whatever you want.
* If you save your shop & reopen it, your SHOP_ID-json-orderimport@inbound.billbee.de E-Mail Address should appear.


### Docs
https://hilfe.billbee.io/article/392-json-e-mail-bestellimport

### Questions?
Contact me on Discord

#### Further information
JK is lazy <3