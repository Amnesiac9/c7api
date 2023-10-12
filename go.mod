module github.com/Amnesiac9/c7api

go 1.20

require github.com/joho/godotenv v1.5.1

retract (
	v1.0.1 // Published by mistake
	v1.0.0 // Published by mistake
)

retract [v1.0.0, v1.0.1]
