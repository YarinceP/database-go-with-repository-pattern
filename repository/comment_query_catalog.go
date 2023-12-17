package repository

// CommentQueries contiene consultas relacionadas con los comentarios.
var CommentQueries = struct {
	Insert        string
	SelectByEmail string
}{
	Insert:        "INSERT INTO comments(email, comment) VALUES(?, ?)",
	SelectByEmail: "SELECT * FROM comments WHERE email = ?",
}
