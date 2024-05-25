//This is a simple To-Do List programme. It was created using Go, HTML and CSS
//To run the programe, open a terminal and type "go run main.go"
//Once it is running, go to your web browser and copy + paste http://localhost:8080
//The To-Do List will appear in a web server running on your local machine at port 8080

package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	todoList []Task
	mu       sync.Mutex
)

type Task struct {
	Description string
	ImageURL    string
	Completed   bool
}

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>To-Do List</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: rgb(65, 32, 107);
            color: white;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 800px;
            margin: auto;
            padding: 20px;
            margin-top: 100px;
            background-color: rgb(42, 19, 72);
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            border-radius: 5px;
        }
        h1 {
            color: white;
            text-align: center;
            margin-bottom: 20px;
        }
        form {
            display: flex;
            margin-bottom: 20px;
        }
        input[type="text"] {
            flex: 1;
            padding: 10px;
            border: 1px solid #ccc;
            border-radius: 5px;
            font-size: 16px;
            margin-right: 10px;
        }
        button[type="submit"] {
            padding: 10px 20px;
            background-color: #4CAF50;
            color: #fff;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
        }
        button[type="submit"]:hover {
            background-color: #45a049;
        }
        ul {
            list-style-type: none;
            padding: 0;
        }
        li {
            margin-bottom: 10px;
            padding: 10px;
            background-color: rgb(184, 153, 225);
            border-radius: 5px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .task-text {
            flex: 1;
            text-decoration: none;
            color: #e9e2e2;
        }
        .delete-btn {
            padding: 5px;
            border: none;
            border-radius: 50%;
            cursor: pointer;
            font-size: 16px;
            background-color: transparent;
            color: white;
            margin-right: 5px;
        }
        .delete-btn:hover {
            background-color: rgba(255, 255, 255, 0.3);
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>To-Do List</h1>
        <form action="/add" method="post">
            <input type="text" name="task" placeholder="Enter task..." required>
            <button type="submit">Add Task</button>
        </form>
        <form action="/" method="get">
            <input type="text" name="search" placeholder="Search...">
            <button type="submit">Search</button>
        </form>
        <ul>
            {{range $index, $task := .}}
            <li>
                <span class="task-text">{{$task.Description}}</span>
                <button class="delete-btn" onclick="deleteTask({{$index}})">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                        <path d="M0 0h24v24H0z" fill="none"/>
                        <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"/>
                    </svg>
                </button>
            </li>
            {{end}}
        </ul>
    </div>
    <script>
        function deleteTask(index) {
            fetch('/delete?index=' + index, { method: 'GET' })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                // Task deleted successfully, update UI if necessary
            })
            .catch(error => {
                console.error('There was a problem with your fetch operation:', error);
            });
        }
    </script>
</body>
</html>
`))

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Check if the search parameter is in URL
	search := r.URL.Query().Get("search")
	if search != "" {
		filteredTodoList := searchTasks(search)
		tmpl.Execute(w, filteredTodoList)
		return
	}

	// Render the template with To-Do List
	tmpl.Execute(w, todoList)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		mu.Lock()
		defer mu.Unlock()
		todoList = append(todoList, Task{Description: r.FormValue("task"), ImageURL: "", Completed: false})
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		mu.Lock()
		defer mu.Unlock()
		index, err := strconv.Atoi(r.URL.Query().Get("index"))
		if err == nil && index >= 0 && index < len(todoList) {
			todoList = append(todoList[:index], todoList[index+1:]...)
			// Respond with success message
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})
			return
		}
	}
	//Respond with error message
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
}

func searchTasks(search string) []Task {
	var result []Task
	for _, task := range todoList {
		if strings.Contains(strings.ToLower(task.Description), strings.ToLower(search)) {
			result = append(result, task)
		}
	}
	return result
}
