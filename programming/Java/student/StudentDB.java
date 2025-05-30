package info.kgeorgiy.ja.kubesh.student;

import info.kgeorgiy.java.advanced.student.GroupName;
import info.kgeorgiy.java.advanced.student.StudentQuery;
import info.kgeorgiy.java.advanced.student.Student;

import java.util.*;
import java.util.function.BinaryOperator;
import java.util.function.Function;
import java.util.function.Predicate;
import java.util.stream.Collectors;

public class StudentDB implements StudentQuery {

    private <T> List<T> getStudentAttributes(List<Student> students, Function<Student, T> mapper) {
        return students.stream()
                .map(mapper)
                .collect(Collectors.toList());
    }

    private List<Student> filterAndSortStudents(Collection<Student> students, Predicate<Student> filter) {
        return students.stream()
                .filter(filter)
                .sorted(Comparator
                        .comparing(Student::firstName)
                        .thenComparing(Student::lastName)
                        .thenComparingInt(Student::id))
                .collect(Collectors.toList());
    }

    @Override
    public List<String> getFirstNames(List<Student> students) {
        return getStudentAttributes(students, Student::firstName);
    }

    @Override
    public List<String> getLastNames(List<Student> students) {
        return getStudentAttributes(students, Student::lastName);
    }

    @Override
    public List<GroupName> getGroupNames(List<Student> students) {
        return getStudentAttributes(students, Student::groupName);
    }

    @Override
    public List<String> getFullNames(List<Student> students) {
        return getStudentAttributes(
                students,
                student -> student.firstName() + " " + student.lastName()
        );
    }

    @Override
    public Set<String> getDistinctFirstNames(List<Student> students) {
        return students.stream()
                .map(Student::firstName)
                .collect(Collectors.toCollection(TreeSet::new));
    }

    @Override
    public String getMaxStudentFirstName(List<Student> students) {
        return students.stream()
                .max(Comparator.comparingInt(Student::id))
                .map(Student::firstName)
                .orElse("");
    }

    @Override
    public List<Student> sortStudentsById(Collection<Student> students) {
        return students.stream()
                .sorted(Comparator.comparingInt(Student::id))
                .collect(Collectors.toList());
    }

    @Override
    public List<Student> sortStudentsByName(Collection<Student> students) {
        return filterAndSortStudents(students, student -> true);
    }

    @Override
    public List<Student> findStudentsByFirstName(Collection<Student> students, String name) {
        return filterAndSortStudents(students, student -> student.firstName().equals(name));
    }

    @Override
    public List<Student> findStudentsByLastName(Collection<Student> students, String name) {
        return filterAndSortStudents(students, student -> student.lastName().equals(name));
    }

    @Override
    public List<Student> findStudentsByGroup(Collection<Student> students, GroupName group) {
        return filterAndSortStudents(students, student -> student.groupName().equals(group));
    }

    @Override
    public Map<String, String> findStudentNamesByGroup(Collection<Student> students, GroupName group) {
        return students.stream()
                .filter(student -> student.groupName().equals(group))
                .collect(Collectors.toMap(
                        Student::lastName,
                        Student::firstName,
                        BinaryOperator.minBy(String::compareToIgnoreCase)
                ));
    }
}
