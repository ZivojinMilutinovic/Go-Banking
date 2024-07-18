We need Docker desktop installed to run this.

In the console(cmd, iterm etc.) docker-compose up --build to run this setup
If you want to stop this setup run  docker-compose down.

Routes for testing:

-- create user
curl -X POST http://localhost:8080/create-user -H "Content-Type: application/json" -d "{\"email\":\"user1@example.com\"}" 

-- transfer money 
curl -X PATCH http://localhost:8081/transfer -H "Content-Type: application/json" -d "{\"from_user_id\": 1, \"to_user_id\": 2, \"amount_to_transfer\": 50.0}"
-- update user balance
curl -X PATCH http://localhost:8081/users/4/balance -H "Content-Type: application/json" -d "{\"amount\": 50.0}"

-- gets user balance
curl "http://localhost:8080/user-balance?email=user1@example.com"