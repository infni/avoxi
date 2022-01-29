curl -X POST http://127.0.0.1:9080/ipauthorize \
   -H 'Content-Type: application/json' \
   -d '{"Ip":"152.216.7.110", "CountryNames":["United States", "Norway"]}'