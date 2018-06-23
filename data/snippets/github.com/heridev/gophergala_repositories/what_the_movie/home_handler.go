package main

import (
	`fmt`
	`net/http`
)

func homeHandler(writer http.ResponseWriter, req *http.Request) {
	resp := `<html>
<body>
	<table>
		<tr>
			<form action="/search" method="POST" enctype="multipart/form-data">
				<td>Search</td>
				<td><input type="file" name="photo"></td>
				<td><input type="submit"></td>
			</form>
		</tr>
	</table>
</body>
</html>`

	fmt.Fprintf(writer, `%s`, resp)
}
