package grpcmarshalsample

import (
	"app/pb"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
)

func GrpcMarshalSample() {
	employee := &pb.Employee{
		Id:         1,
		Name:       "愛はあるんか",
		Email:      "愛はあるんか@love.com",
		Ocupation:  pb.Ocupation_ENGINEER,
		PhonNumber: []string{"080-love-exist", "090-that-love-exist"},
		// Project:    map[string]*pb.Project{"ProjectX": &pb.Project{"": }},
		Profile: &pb.Employee_Text{
			Text: "My name is 愛はあるんか",
		},
		Birth: &pb.Date{Year: 2023},
	}

	// binData, _ := proto.Marshal(employee)

	// var fileName string = "em.bin"
	// if err := os.WriteFile(fileName, binData, 0666); err != nil {
	// 	log.Fatalln("eeror: ", err)
	// }

	// fmt.Println(employee.GetBirth())

	// readFile, er := os.ReadFile(fileName)

	// fmt.Println("er: ", er)

	// // fmt.Println(string(readFile))

	// readData := &pb.Employee{}
	// err := proto.Unmarshal(readFile, readData)
	// fmt.Println("error: ", err)
	// fmt.Println(readData)

	m := jsonpb.Marshaler{}

	out, _ := m.MarshalToString(employee)
	fmt.Println(out)

	readEmployee := &pb.Employee{}
	jsonpb.UnmarshalString(out, readEmployee)

	fmt.Println(readEmployee)
}
