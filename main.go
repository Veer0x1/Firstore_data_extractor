package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	// "firebase.google.com/go/auth"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the credentials.json file")
		return
	}
	if len(os.Args) < 3 {
		fmt.Println("Please provide the collection name")
		return
	}

	fileName := os.Args[1]
	collectionName := os.Args[2]

	opt := option.WithCredentialsFile(fileName)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println("Error initializing app:", err)
		return
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		fmt.Println("Error initializing Firestore client:", err)
		return
	}

	defer client.Close()

	iter := client.Collection(collectionName).Documents(context.Background())
	defer iter.Stop()

	csvFile, err := os.Create(fmt.Sprintf("%s.csv", collectionName))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}

		var row []string
		for _, value := range doc.Data() {
			row = append(row, fmt.Sprintf("%v", value))
		}

		if err := writer.Write(row); err != nil {
			fmt.Println("Error writing record to csv:", err)
			return
		}
	}

	fmt.Println("Data exported successfully")

}
