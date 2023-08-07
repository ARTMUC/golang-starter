package router

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

func generateDocs(controllers []interface{}) {

	filename := "routes/generated_functions.go" // File to save all the functions

	var builder strings.Builder

	builder.WriteString("package routes\n\n")

	for _, ctrl := range controllers {
		controller := ctrl.(Controller)
		pkgPath := reflect.TypeOf(controller).Elem().PkgPath()
		pkgParts := strings.Split(pkgPath, "/")
		pkgName := pkgParts[len(pkgParts)-1]
		fmt.Println("pkg name", pkgName)

		baseUrl := controller.MainPath()

		//listFields(ctrl)
		pn := pkgName
		postFunctionName := pn + "PostGenerated"
		listFunctionName := pn + "ListGenerated"
		getFunctionName := pn + "GetGenerated"
		updateFunctionName := pn + "UpdateGenerated"
		deleteFunctionName := pn + "DeleteGenerated"
		createDto := pn + ".CreateDto"
		updateDto := pn + ".UpdateDto"
		responseObject := pn + ".Model"

		var omitPost bool
		var omitList bool
		var omitGet bool
		var omitUpdate bool
		var omitDelete bool
		//
		//if omits, ok := value["omit"]; ok && omits != nil {
		//	if omitsArr, ok := omits.([]string); ok && len(omitsArr) > 0 {
		//		for _, omit := range omitsArr {
		//			switch omit {
		//			case "create":
		//				omitPost = true
		//			case "list":
		//				omitList = true
		//			case "get":
		//				omitGet = true
		//			case "update":
		//				omitUpdate = true
		//			case "delete":
		//				omitDelete = true
		//			}
		//		}
		//	}
		//}

		var code string

		if !omitPost {
			code += generatePostFunctions(baseUrl, postFunctionName, responseObject, createDto)
		}

		if !omitUpdate {
			code += generateUpdateFunctions(baseUrl, updateFunctionName, responseObject, updateDto)
		}

		if !omitDelete {
			code += generateDeleteFunctions(baseUrl, deleteFunctionName)
		}

		if !omitGet {
			code += generateGetFunctions(baseUrl, getFunctionName, responseObject)
		}

		if !omitList {
			//var conditions []string
			//if cond, ok := value["conditions"].([]string); ok {
			//	conditions = cond
			//}
			code += generateListFunctions(baseUrl, listFunctionName, responseObject, []string{})
		}

		builder.WriteString(code)
		//}
		//
		err := saveToFile(filename, builder.String())
		if err != nil {
			fmt.Printf("Error saving file '%s': %v\n", filename, err)
		} else {
			fmt.Printf("File '%s' saved successfully.\n", filename)
		}

	}
}

func generatePostFunctions(key, functionName, responseObject string, createDto string) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("// %s godoc\n", functionName))
	builder.WriteString(fmt.Sprintf("// @Summary   Creates entity\n"))
	builder.WriteString(fmt.Sprintf("// @Tags   %s crud   \n", key))
	builder.WriteString(fmt.Sprintf("// @Accept    json\n"))
	builder.WriteString(fmt.Sprintf("// @Produce   json\n"))
	builder.WriteString(fmt.Sprintf("// @Param\t  request body %s true \"Request data\"      \n", createDto))
	builder.WriteString(fmt.Sprintf("// @Response  201 {object}        %s\n", responseObject))
	builder.WriteString(fmt.Sprintf("// @Router    /%s [post]\n", key))
	builder.WriteString(fmt.Sprintf("func %s() {}\n\n", functionName))

	return builder.String()
}

func generateGetFunctions(key, functionName, responseObject string) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("// %s godoc\n", functionName))
	builder.WriteString(fmt.Sprintf("// @Summary   Returns entity\n"))
	builder.WriteString(fmt.Sprintf("// @Tags   %s crud   \n", key))
	builder.WriteString(fmt.Sprintf("// @Accept    json\n"))
	builder.WriteString(fmt.Sprintf("// @Produce   json\n"))
	builder.WriteString(fmt.Sprintf("// @Param\t   id path string true \"id of item\"      \n"))
	builder.WriteString(fmt.Sprintf("// @Response  200 {object}        %s\n", responseObject))
	builder.WriteString(fmt.Sprintf("// @Router    /%s/:id [get]\n", key))
	builder.WriteString(fmt.Sprintf("func %s() {}\n\n", functionName))

	return builder.String()
}

func generateDeleteFunctions(key, functionName string) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("// %s godoc\n", functionName))
	builder.WriteString(fmt.Sprintf("// @Summary   Deletes entity\n"))
	builder.WriteString(fmt.Sprintf("// @Tags   %s crud   \n", key))
	builder.WriteString(fmt.Sprintf("// @Accept    json\n"))
	builder.WriteString(fmt.Sprintf("// @Produce   json\n"))
	builder.WriteString(fmt.Sprintf("// @Param\t   id path string true \"id of item\"      \n"))
	builder.WriteString(fmt.Sprintf("// @Router    /%s/:id [delete]\n", key))
	builder.WriteString(fmt.Sprintf("func %s() {}\n\n", functionName))

	return builder.String()
}

func generateUpdateFunctions(key, functionName, responseObject string, updateDto string) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("// %s godoc\n", functionName))
	builder.WriteString(fmt.Sprintf("// @Summary   Updates entity\n"))
	builder.WriteString(fmt.Sprintf("// @Tags   %s crud   \n", key))
	builder.WriteString(fmt.Sprintf("// @Accept    json\n"))
	builder.WriteString(fmt.Sprintf("// @Produce   json\n"))
	builder.WriteString(fmt.Sprintf("// @Param\t   id path string true \"id of item\"      \n"))
	builder.WriteString(fmt.Sprintf("// @Param\t  request body %s true \"Request data\"      \n", updateDto))
	builder.WriteString(fmt.Sprintf("// @Response  200 {object}        %s\n", responseObject))
	builder.WriteString(fmt.Sprintf("// @Router    /%s/:id [put]\n", key))
	builder.WriteString(fmt.Sprintf("func %s() {}\n\n", functionName))

	return builder.String()
}

func generateListFunctions(key, functionName, responseObject string, conditions []string) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("// %s godoc\n", functionName))
	builder.WriteString(fmt.Sprintf("// @Summary   Lists entities\n"))
	builder.WriteString(fmt.Sprintf("// @Tags   %s crud   \n", key))
	builder.WriteString(fmt.Sprintf("// @Accept    json\n"))
	builder.WriteString(fmt.Sprintf("// @Produce   json\n"))
	builder.WriteString(fmt.Sprintf("// @Param     s                   query string false  \"{'$and': [ {'title': { '$cont':'cul' } } ]}\"\"\n"))
	builder.WriteString(fmt.Sprintf("// @Param     fields               query string false  \"fields to select eg: name,age\"\n"))
	builder.WriteString(fmt.Sprintf("// @Param     page             query string false  \"page of pagination\"\n"))
	builder.WriteString(fmt.Sprintf("// @Param     limit            query string false  \"limit of pagination\"\n"))
	builder.WriteString(fmt.Sprintf("// @Param     join            query string false  \"join relations eg: category, parent\"\n"))
	builder.WriteString(fmt.Sprintf("// @Param     filter            query string false  \"filters eg: name||$eq||ad price||$gte||200\"\n"))
	builder.WriteString(fmt.Sprintf("// @Param     sort            query string false  \"filters eg: created_at,desc title,asc\"\n"))
	for _, condition := range conditions {
		builder.WriteString(fmt.Sprintf("// @Param     %s            query string false  \"%s\"\n", condition, condition))
	}
	builder.WriteString(fmt.Sprintf("// @Response  200 {object}        pagination.Output{rows=[]%s} \n", responseObject))
	builder.WriteString(fmt.Sprintf("// @Router    /campaign/crud/:campaign/%s [get]\n", key))
	builder.WriteString(fmt.Sprintf("func %s() {}\n\n", functionName))

	return builder.String()
}

func saveToFile(filename, code string) error {
	return ioutil.WriteFile(filename, []byte(code), 0644)
}

//func listFields(obj interface{}) {
//	t := reflect.TypeOf(obj).Elem()
//	if t.Kind() != reflect.Struct {
//		fmt.Println("Expected a struct type.")
//		return
//	}
//
//	for i := 0; i < t.NumField(); i++ {
//		field := t.Field(i)
//		fmt.Println("---------------")
//		fmt.Println("Field:", field.Name)
//		fmt.Println("Type:", field.Type)
//		fmt.Println("Tag:", field.Tag)
//		fmt.Println("---------------")
//
//		// Check if the field is a struct
//		if field.Type.Kind() == reflect.Ptr {
//			if field.Type.Elem().Kind() == reflect.Struct {
//				// Recursively list fields of nested struct
//				listFields(reflect.New(field.Type).Elem().Interface())
//			}
//		}
//
//	}
//}
