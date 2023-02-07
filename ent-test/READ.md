

C:\atlas\atlas-windows-amd64-latest.exe  migrate apply --dir "file://ent/migrate/migrations" --url postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable

C:\atlas\atlas-windows-amd64-latest.exe  migrate apply --dir "file://migrations/db" --url postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable

C:\atlas\atlas-windows-amd64-latest.exe  migrate hash --dir "file://migrations/db" 