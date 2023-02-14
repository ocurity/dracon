package example.gosec

default allow := false

allow =true {
    print(input)
    check_severity
}

check_severity {
    input.severity == "SEVERITY_LOW"
}

check_severity {
    input.severity == "SEVERITY_HIGH"
}
check_severity {
    input.severity == "SEVERITY_MEDIUM"
}
