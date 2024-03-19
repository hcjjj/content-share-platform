--wrk -t1 -d1s -c2 -s ./scripts/wrk/profile.lua http://localhost:8080/users/profile

wrk.method="GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["User-Agent"] = "PostmanRuntime/7.32.3"
-- 记得修改这个，你在登录页面登录一下，然后复制一个过来这里
wrk.headers["Authorization"]="Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6NiwiVXNlckFnZW50IjoiUG9zdG1hblJ1bnRpbWUvNy4zMi4zIiwiZXhwIjoxNjkwMjczNjUwfQ.qmZ2jwT-JxDy4uGpuKJLSudEDpoxC1FDOe_KciNZbO8"