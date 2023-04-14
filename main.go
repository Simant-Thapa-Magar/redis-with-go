package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Student struct {
	Name  string
	Total int
}

func main() {
	ctx := context.Background()
	client := NewClient(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASS"), 0)

	if err := client.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection established successfully")

	if err := client.Set(ctx, "my_key", "my_value", 0); err != nil {
		log.Fatal(err)
	}

	value, err := client.Get(ctx, "my_key")

	if err != nil {
		fmt.Println("Error getting my_key")
	}

	fmt.Println("my_key is ", value)

	students := []Student{
		{
			Name:  "Student1",
			Total: 74,
		},
		{
			Name:  "Student2",
			Total: 46,
		},
		{
			Name:  "Student3",
			Total: 91,
		},
	}

	dataInByte, err := json.Marshal(students)

	if err != nil {
		fmt.Println("error marshalling")
		log.Fatal(err)
	}

	if err := client.Set(ctx, "students", dataInByte, 0); err != nil {
		log.Fatal(err)
	}

	asbyte, _ := client.Get(ctx, "students")

	unmarshalled := []Student{}
	json.Unmarshal([]byte(asbyte), &unmarshalled)
	fmt.Println("Students are ", unmarshalled)

	client.SetStudentsWithScore(ctx, students)

	rankedStudents, err := client.GetStudentsByRank(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Student with rank based on total")
	for rank, student := range rankedStudents {
		fmt.Println("Rank: ", rank+1, " Student ", student)
	}
}
