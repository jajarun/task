package main

import (
	"fmt"
	"reflect"
	"task/models"
)

type user struct {
	Name string
}

func (u *user) test() {
	fmt.Println(u.Name)
}
func test(t interface{}) {
	ptrValue := reflect.ValueOf(t)
	if ptrValue.Kind() != reflect.Ptr || ptrValue.Elem().Kind() != reflect.Slice {
		fmt.Println("Invalid argument")
		return
	}

	sliceValue := ptrValue.Elem()
	for i := 0; i < sliceValue.Len(); i++ {
		elemValue := sliceValue.Index(i)
		if elemValue.Kind() == reflect.Ptr {
			elemValue = elemValue.Elem()
		}
		if elemValue.Kind() != reflect.Struct {
			fmt.Println("Invalid slice element")
			continue
		}
		if !elemValue.Addr().MethodByName("test").IsValid() {
			fmt.Println("Method not found")
			continue
		}
		args := []reflect.Value{}
		elemValue.Addr().MethodByName("test").Call(args)
	}
}

func test1() {
	//users := []user{{Name: "111"}, {Name: "222"}}
	//test(&users)
	user := models.User{}
	models.FindByPrimaryKey(&user, "18")
	//fmt.Println(user)
	user.UserName = "jajatest"
	models.Save(&user)
	//fmt.Println(user)

	//users := []models.User{}
	//models.Find(&users)
	//fmt.Println(users)

}
