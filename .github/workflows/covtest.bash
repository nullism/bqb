last_cov=$(go tool cover -func=coverage.out | tail -1 | awk '{ print $3 }' | sed s/%//g)
go test -coverprofile=coverage.out
this_cov=$(go tool cover -func=coverage.out | tail -1 | awk '{ print $3 }' | sed s/%//g)

if [ "${last_cov}" -gt "${this_cov}" ]; then
    echo "Coverage has dropped from ${last_cov} to ${this_cov}"
    exit 1
fi