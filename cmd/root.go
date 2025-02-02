/*
Copyright © 2022 Bora Kasmer <bora@borakasmer.com>

*/
//go get github.com/olekukonko/tablewriter
package cmd

import (
	"github.com/borakasmer/fuel/Model"
	"github.com/borakasmer/fuel/parser"
	"github.com/olekukonko/tablewriter"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fuel",
	Short: "Bu Cli Tool'u ile, PO'dan güncel petrol fiyatları İstanbul, Ankara ve İzmir için çekilir.",
	Long: `fuel

**** Herhangi bir tanımlı parametre yoktur.
------------------
Petrol Ofisi sitesi	üzerindeki fiyat listesi, anlık olarak Parse Edilerek ekrana, İl bazında basılır.
Örnek kullanım:
."fuel"
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		getFuel()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fuel.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getFuel() {
	tableHeaders := make([]string, 0)
	tableRows := make([][]string, 0)

	tableHeaders = append(tableHeaders, "Yakıt Tipi")

	var ExchangeList = make(map[string]string, 3)
	ExchangeList["Istanbul"] = parser.IstanbulUrl
	ExchangeList["Angara"] = parser.AnkaraUrl
	ExchangeList["Izmir"] = parser.IzmirUrl

	resultList := make([]Model.FuelPrice, 0)

	for key, url := range ExchangeList {
		petrol, diesel, lpg := parser.ParseWeb(url)

		resultList = append(resultList,
			Model.FuelPrice{City: key,
				Diesel:      diesel.Slice(),
				Petrol:      petrol.Slice(),
				Lpg:         lpg.Slice(),
				CurrentDate: time.Now().Format("2006.01.02 15:04:05"),
			})
	}
	petrolRow := make([]string, 0)
	petrolRow = append(petrolRow, "Benzin =>")
	diesel := make([]string, 0)
	diesel = append(diesel, "Mazot =>")
	lpg := make([]string, 0)
	lpg = append(lpg, "Lpg =>")
	for _, item := range resultList {
		tableHeaders = append(tableHeaders, item.City)
		petrolRow = append(petrolRow, item.Petrol+"₺")
		diesel = append(diesel, item.Diesel+"₺")
		lpg = append(lpg, item.Lpg+"₺")
	}

	tableHeaders = append(tableHeaders, "Tarih")
	petrolRow = append(petrolRow, time.Now().Format("2006.01.02 15:04:05"))
	diesel = append(diesel, time.Now().Format("2006.01.02 15:04:05"))
	lpg = append(lpg, time.Now().Format("2006.01.02 15:04:05"))

	tableRows = append(tableRows, petrolRow)
	tableRows = append(tableRows, diesel)
	tableRows = append(tableRows, lpg)

	//Create Header of Table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tableHeaders)
	//---------------------

	table.SetCaption(true, "Petrol Ofisinde, Şehir Bazlı Yakıt Fiyatları\n ®coderbora => www.borakasmer.com")
	table.AppendBulk(tableRows)
	//Set Table Color
	if !isWindows() { //Windows için Renkli Tablo başlıkları gözükmüyor...

		table.SetHeaderColor(tablewriter.Colors{
			tablewriter.Bold, tablewriter.BgMagentaColor},
			tablewriter.Colors{
				tablewriter.Bold, tablewriter.BgGreenColor},
			tablewriter.Colors{
				tablewriter.Bold, tablewriter.BgYellowColor},
			tablewriter.Colors{
				tablewriter.Bold, tablewriter.BgBlueColor},
			tablewriter.Colors{
				tablewriter.Bold, tablewriter.BgRedColor})
	}
	table.Render()
}

func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
