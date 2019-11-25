package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	goengage "github.com/salsalabs/goengage/pkg"
)

func show(a []goengage.Supporter) {
	for i, s := range a {
		fmt.Printf("%2d %v\n", i+1, text(s))
	}
}
func text(s goengage.Supporter) string {
	return fmt.Sprintf("%v %v %v",
		s.FirstName,
		s.LastName,
		s.Suffix)
}
func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&goengage.Supporter{})
	db.AutoMigrate(&goengage.Contact{})
	db.AutoMigrate(&goengage.CustomFieldValue{})

	os.Exit(0)

	// Create
	db.Create(&goengage.Supporter{FirstName: "Bubba", LastName: "Seconal", Suffix: "bubba@seconal.biz"})
	db.Create(&goengage.Supporter{FirstName: "Cissy", LastName: "Seconal", Suffix: "cissy@seconal.biz"})
	db.Create(&goengage.Supporter{FirstName: "PawPaw", LastName: "Seconal", Suffix: "pawpaw@seconal.biz"})
	db.Create(&goengage.Supporter{FirstName: "MeeMaw", LastName: "Seconal", Suffix: "meemaw@seconal.biz"})
	db.Create(&goengage.Supporter{FirstName: "Pop", LastName: "Seconal", Suffix: "pop@seconal.biz"})
	db.Create(&goengage.Supporter{FirstName: "Moms", LastName: "Seconal", Suffix: "moms@seconal.biz"})

	// Read
	var s goengage.Supporter
	db.First(&s, "Suffix = ?", "moms@seconal.biz")
	fmt.Printf("First with suffix: %v\n", text(s))

	db.Model(&s).Update("LastName", "Blizzard")
	db.First(&s, "Suffix = ?", "moms@seconal.biz")
	fmt.Printf("Confirm update: %v\n", text(s))

	var a []goengage.Supporter

	fmt.Println("Read all")
	db.Find(&a)
	show(a)

	fmt.Println("Scrub")
	for i, s := range a {
		db.Delete(&s)
		fmt.Printf("%2d %s deleted\n", i+1, text(s))
	}
}
