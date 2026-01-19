mysql://user:password@tcp(host:port)/dbname?query

migrate create -ext sql -dir db/migrations -seq create_order_menu_items_table #create a new migration file

migrate -database "mysql://root:root@tcp(127.0.0.1:3306)/api_coffee_shop" -path db/migrations up #run all up migrations
migrate -database "mysql://root:root@tcp(127.0.0.1:3306)/api_coffee_shop" -path db/migrations down #run all down migrations
migrate -database "mysql://root:root@tcp(127.0.0.1:3306)/api_coffee_shop" -path db/migrations force N #set version to N without running migrations

migrate -database "mysql://adiva:jaringan123@tcp(43.156.16.6:3306)/api_coffee_shop" -path db/migrations up