package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"webserver/internal/entities"
	"webserver/internal/services"
	"webserver/pkg/utils"

	"github.com/gorilla/mux"
)

type BatHandler struct {
	batService services.BatService
}

func NewBatHandler(batService services.BatService) *BatHandler {
	return &BatHandler{
		batService: batService,
	}
}

func (h *BatHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", h.Info).Methods("GET")
	router.HandleFunc("/{action}/{param}", h.HandleBatRequest).Methods("GET")
	router.HandleFunc("/api/execute", h.HandleJSONRequest).Methods("POST")
	router.HandleFunc("/health", h.HealthCheck).Methods("GET")
	router.HandleFunc("/stats", h.GetStats).Methods("GET")
}

func (h *BatHandler) Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Batch File Executor Service</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        
        .header {
            background: white;
            padding: 40px;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            text-align: center;
            margin-bottom: 30px;
        }
        
        .header h1 {
            color: #4a5568;
            font-size: 2.5rem;
            margin-bottom: 10px;
        }
        
        .header p {
            color: #718096;
            font-size: 1.2rem;
        }
        
        .content {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 30px;
            margin-bottom: 30px;
        }
        
        .card {
            background: white;
            padding: 30px;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
        }
        
        .card h2 {
            color: #4a5568;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 3px solid #667eea;
        }
        
        .example {
            background: #f7fafc;
            padding: 20px;
            border-radius: 10px;
            margin: 15px 0;
            border-left: 4px solid #667eea;
        }
        
        .code {
            background: #2d3748;
            color: #e2e8f0;
            padding: 15px;
            border-radius: 8px;
            font-family: 'Courier New', monospace;
            margin: 10px 0;
            overflow-x: auto;
        }
        
        .endpoint {
            background: #edf2f7;
            padding: 15px;
            border-radius: 8px;
            margin: 10px 0;
            border: 1px solid #e2e8f0;
        }
        
        .method {
            display: inline-block;
            background: #48bb78;
            color: white;
            padding: 4px 12px;
            border-radius: 4px;
            font-weight: bold;
            margin-right: 10px;
        }
        
        .method.get { background: #48bb78; }
        .method.post { background: #4299e1; }
        
        .footer {
            text-align: center;
            color: white;
            margin-top: 40px;
        }
        
        @media (max-width: 768px) {
            .content {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🛠️ Batch File Executor Service</h1>
            <p>Веб-сервис для управления базами данных через HTTP запросы</p>
        </div>
        
        <div class="content">
            <div class="card">
                <h2>🚀 Быстрый старт</h2>
                <p>Используйте следующие endpoints для создания и обновления баз данных:</p>
                
                <div class="endpoint">
                    <span class="method get">GET</span>
                    <strong>/create/&lt;имя_базы&gt;</strong>
                </div>
                <div class="example">
                    <strong>Пример:</strong>
                    <div class="code">http://localhost:5005/create/my_database</div>
                    <p>Создает новую базу данных с указанным именем</p>
                </div>
                
                <div class="endpoint">
                    <span class="method get">GET</span>
                    <strong>/update/&lt;имя_базы&gt;</strong>
                </div>
                <div class="example">
                    <strong>Пример:</strong>
                    <div class="code">http://localhost:5005/update/my_database</div>
                    <p>Обновляет существующую базу данных</p>
                </div>
            </div>
            
            <div class="card">
                <h2>📡 JSON API</h2>
                <p>Для интеграции с приложениями используйте JSON API:</p>
                
                <div class="endpoint">
                    <span class="method post">POST</span>
                    <strong>/api/execute</strong>
                </div>
                
                <div class="example">
                    <strong>Запрос:</strong>
                    <div class="code">{
    "action": "create",
    "param": "my_database"
}</div>
                    
                    <strong>Ответ:</strong>
                    <div class="code">{
    "success": true,
    "message": "База данных 'my_database' успешно создана",
    "output": "Вывод выполнения скрипта...",
    "timestamp": "2024-01-15T10:30:00Z"
}</div>
                </div>
            </div>
        </div>
        
        <div class="footer">
            <p>Batch File Executor Service &copy; 2024 | Версия 1.0.0</p>
        </div>
    </div>
</body>
</html>`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// Остальные методы остаются без изменений
func (h *BatHandler) HandleBatRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	param := vars["param"]

	// Валидация параметра
	if !utils.IsValidParameter(param) {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid parameter format")
		return
	}

	// Обрабатываем запрос
	response, err := h.batService.ProcessRequest(r.Context(), action, param)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.sendSuccessResponse(w, response)
}

func (h *BatHandler) HandleJSONRequest(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Action string `json:"action"`
		Param  string `json:"param"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Валидация
	if !utils.IsValidParameter(request.Param) {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid parameter format")
		return
	}

	// Обрабатываем запрос
	response, err := h.batService.ProcessRequest(r.Context(), request.Action, request.Param)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.sendSuccessResponse(w, response)
}

func (h *BatHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "bat-executor",
		"version":   "1.0.0",
	})
}

func (h *BatHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := h.batService.GetStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"stats":     stats,
		"timestamp": time.Now(),
	})
}

func (h *BatHandler) handleServiceError(w http.ResponseWriter, err error) {
	h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
}

func (h *BatHandler) sendSuccessResponse(w http.ResponseWriter, response *entities.BatResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *BatHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	errorResponse := map[string]interface{}{
		"success":   false,
		"error":     message,
		"timestamp": time.Now(),
	}

	json.NewEncoder(w).Encode(errorResponse)
}
