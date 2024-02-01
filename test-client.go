package main

import (
	"context"
	"fmt"
	"github.com/minio/madmin-go/v3"
	"net/url"
	"reflect"
	"time"
	"unsafe"
)

func main() {
	endpoint := "https://acme-0-0.lab.min.dev:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, "UserID", 123)

	printContext(ctx, true)

	admClient, err := getAdminClient(endpoint, accessKeyID, secretAccessKey)
	if err != nil {
		fmt.Println(err.Error())
	}
	print(admClient)
	if err = admClient.AddUser(ctx, "newuser", "newstrongpassword"); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("user added")
}

func getAdminClient(endpoint, accessKey, secretKey string) (*madmin.AdminClient, error) {
	epURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	client, err := madmin.New(epURL.Host, accessKey, secretKey, epURL.Scheme == "https")
	if err != nil {
		return nil, err
	}
	return client, nil
}

func printContext(ctx interface{}, inner bool) {
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	if !inner {
		fmt.Printf("\nFields for %s.%s\n", contextKeys.PkgPath(), contextKeys.Name())
	}

	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()
			reflectField := contextKeys.Field(i)
			if reflectField.Name == "Context" {
				printContext(reflectValue.Interface(), true)
			} else {
				fmt.Printf("Field: %v, Value: %v\n", reflectField.Name, reflectValue.Interface())
			}
		}
	} else {
		fmt.Printf("Context is empty\n")
	}
}
