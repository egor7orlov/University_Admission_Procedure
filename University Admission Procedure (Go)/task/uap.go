package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

const studentsListFilename = "applicants.txt"
const applicationWavesAmount = 3

func main() {
	departmentCapacity := GetIntFromInputOrFail()
	students := CollectStudentsFromFile(studentsListFilename)
	dptsEnrolment := CreateDepartmentsEnrolment(departmentCapacity, students)

	dptsEnrolment.ProceedEnrolment(applicationWavesAmount)
	dptsEnrolment.WriteAcceptedStudentsToFiles()
}

// Common

func GetIntFromInputOrFail() int {
	var number int

	_, err := fmt.Scan(&number)

	if err != nil {
		log.Fatal(err)
	}

	return number
}

func ParseFloat64OrFail(s string) float64 {
	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

// Student

type Student struct {
	FirstName      string
	LastName       string
	DptsGrades     map[DepartmentName]float64
	DptsByPriority map[int]DepartmentName
}

type FullName = string

func (student *Student) FullName() FullName {
	return fmt.Sprintf("%s %s", student.FirstName, student.LastName)
}

func CollectStudentsFromFile(fileName string) []Student {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	students := make([]Student, 0, 100)

	for scanner.Scan() {
		studentDataStr := scanner.Text()
		student := ParseStudentFromString(studentDataStr)
		students = append(students, student)
	}

	return students
}

func ParseStudentFromString(s string) Student {
	studentData := strings.Split(s, " ")
	firstName := studentData[0]
	lastName := studentData[1]
	gradePhysics := ParseFloat64OrFail(studentData[2])
	gradeChemistry := ParseFloat64OrFail(studentData[3])
	gradeMath := ParseFloat64OrFail(studentData[4])
	gradeCS := ParseFloat64OrFail(studentData[5])
	gradeAdmission := ParseFloat64OrFail(studentData[6])
	firstPriorityDpt := studentData[7]
	secondPriorityDpt := studentData[8]
	thirdPriorityDpt := studentData[9]

	return Student{
		FirstName: firstName,
		LastName:  lastName,
		DptsGrades: map[DepartmentName]float64{
			DepartmentPhysics:     math.Max((gradePhysics+gradeMath)/2, gradeAdmission),
			DepartmentChemistry:   math.Max(gradeChemistry, gradeAdmission),
			DepartmentMathematics: math.Max(gradeMath, gradeAdmission),
			DepartmentEngineering: math.Max((gradeCS+gradeMath)/2, gradeAdmission),
			DepartmentBiotech:     math.Max((gradeChemistry+gradePhysics)/2, gradeAdmission),
		},
		DptsByPriority: map[int]DepartmentName{
			1: firstPriorityDpt,
			2: secondPriorityDpt,
			3: thirdPriorityDpt,
		},
	}
}

func SortStudentsByDptGrade(students *[]Student, department DepartmentName) {
	sort.Slice(*students, func(i, j int) bool {
		student1 := (*students)[i]
		student2 := (*students)[j]

		if student1.DptsGrades[department] == student2.DptsGrades[department] {
			return student1.FullName() < student2.FullName()
		}

		return student1.DptsGrades[department] > student2.DptsGrades[department]
	})
}

func GetMapOfStudentsFromSlice(students []Student) map[FullName]Student {
	studentsMap := make(map[FullName]Student)

	for _, student := range students {
		studentsMap[student.FullName()] = student
	}

	return studentsMap
}

func GetSliceOfStudentsFromMap(studentsMap map[FullName]Student) []Student {
	students := make([]Student, 0, len(studentsMap))

	for _, student := range studentsMap {
		students = append(students, student)
	}

	return students
}

func GetStudentsByDepartment(students []Student, department DepartmentName, priority int) []Student {
	result := make([]Student, 0, len(students)/2)

	for _, student := range students {
		if student.DptsByPriority[priority] == department {
			result = append(result, student)
		}
	}

	SortStudentsByDptGrade(&result, department)

	return result
}

// Department

const (
	DepartmentMathematics = "Mathematics"
	DepartmentPhysics     = "Physics"
	DepartmentBiotech     = "Biotech"
	DepartmentChemistry   = "Chemistry"
	DepartmentEngineering = "Engineering"
)

type DepartmentName = string

type DepartmentsEnrolment struct {
	DepartmentsCapacity int
	DepartmentsStudents map[DepartmentName][]Student
	AllStudents         []Student
}

func CreateDepartmentsEnrolment(departmentCapacity int, allStudents []Student) DepartmentsEnrolment {
	return DepartmentsEnrolment{
		AllStudents:         allStudents,
		DepartmentsCapacity: departmentCapacity,
		DepartmentsStudents: map[DepartmentName][]Student{
			DepartmentMathematics: make([]Student, 0, departmentCapacity),
			DepartmentPhysics:     make([]Student, 0, departmentCapacity),
			DepartmentBiotech:     make([]Student, 0, departmentCapacity),
			DepartmentChemistry:   make([]Student, 0, departmentCapacity),
			DepartmentEngineering: make([]Student, 0, departmentCapacity),
		},
	}
}

func (de *DepartmentsEnrolment) ProceedEnrolment(wavesAmount int) {
	for i := 1; i <= wavesAmount; i++ {
		de.ProceedEnrolmentWave(i)
	}
}

func (de *DepartmentsEnrolment) ProceedEnrolmentWave(priority int) {
	dptsNames := [5]DepartmentName{
		DepartmentMathematics,
		DepartmentEngineering,
		DepartmentPhysics,
		DepartmentBiotech,
		DepartmentChemistry,
	}

	for _, dpt := range dptsNames {
		de.GatherSuccessfulApplicantsForDepartment(dpt, priority)
	}
}

func (de *DepartmentsEnrolment) GatherSuccessfulApplicantsForDepartment(
	department DepartmentName,
	priority int,
) {
	currentDptStudentsAmt := len(de.DepartmentsStudents[department])

	if currentDptStudentsAmt >= de.DepartmentsCapacity {
		return
	}

	studentsByPrioritizedDpt := GetStudentsByDepartment(de.AllStudents, department, priority)

	if len(studentsByPrioritizedDpt) <= 0 {
		return
	}

	var topStudents []Student
	acceptableStudentsAmt := de.DepartmentsCapacity - currentDptStudentsAmt

	if acceptableStudentsAmt > len(studentsByPrioritizedDpt) {
		topStudents = studentsByPrioritizedDpt
	} else {
		topStudents = studentsByPrioritizedDpt[:acceptableStudentsAmt]
	}

	if len(topStudents) <= 0 {
		return
	}

	de.DepartmentsStudents[department] = append(de.DepartmentsStudents[department], topStudents...)
	studentsMap := GetMapOfStudentsFromSlice(de.AllStudents)

	for _, student := range topStudents {
		delete(studentsMap, student.FullName())
	}

	de.AllStudents = GetSliceOfStudentsFromMap(studentsMap)
}

func (de *DepartmentsEnrolment) GetDptStudentsLines(dptName DepartmentName) string {
	var result string

	dptStudents, ok := de.DepartmentsStudents[dptName]

	SortStudentsByDptGrade(&dptStudents, dptName)

	if !ok {
		log.Fatal("Unknown department name")
	}

	for _, student := range dptStudents {
		result = fmt.Sprintf("%s\n%s %.2f", result, student.FullName(), student.DptsGrades[dptName])
	}

	return strings.TrimSpace(result)
}

func (de *DepartmentsEnrolment) WriteAcceptedStudentsToFiles() {
	dptsNames := []DepartmentName{
		DepartmentMathematics,
		DepartmentEngineering,
		DepartmentPhysics,
		DepartmentBiotech,
		DepartmentChemistry,
	}

	for _, dptName := range dptsNames {
		fileName := strings.ToLower(dptName) + ".txt"
		studentsLines := de.GetDptStudentsLines(dptName)
		if err := os.WriteFile(fileName, []byte(studentsLines), 0644); err != nil {
			log.Fatal(err)
		}
	}
}
