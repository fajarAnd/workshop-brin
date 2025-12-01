package handlers

import (
	"net/http"

	"github.com/fajarAnd/workshop-brin/wa-service/internal/app/services"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type QRHandler interface {
	GetQRCode(c *gin.Context)
	GetConnectionStatus(c *gin.Context)
	ShowQRPage(c *gin.Context)
	GetQRImage(c *gin.Context)
}

type qrHandler struct {
	whatsappService services.WhatsAppService
}

func NewQRHandler(whatsappService services.WhatsAppService) QRHandler {
	return &qrHandler{
		whatsappService: whatsappService,
	}
}

func (h *qrHandler) GetQRCode(c *gin.Context) {
	qrCode, err := h.whatsappService.GetQRCode()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "QR code not available. Please ensure WhatsApp service is running and not connected",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"qr_code":     qrCode,
		"message":     "Scan this QR code with your WhatsApp to connect the bot",
		"instruction": "1. Open WhatsApp on your phone\n2. Go to Settings > Linked Devices\n3. Tap 'Link a Device'\n4. Scan this QR code",
	})
}

func (h *qrHandler) GetConnectionStatus(c *gin.Context) {
	isConnected := h.whatsappService.IsConnected()

	status := "disconnected"
	message := "WhatsApp bot is not connected. Please scan QR code to connect."

	if isConnected {
		status = "connected"
		message = "WhatsApp bot is connected and ready to receive messages."
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"connected": isConnected,
		"status":    status,
		"message":   message,
	})
}

func (h *qrHandler) ShowQRPage(c *gin.Context) {
	isConnected := h.whatsappService.IsConnected()

	if isConnected {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
<!DOCTYPE html>
<html>
<head>
    <title>WhatsApp Bot - Connected</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { 
            font-family: Arial, sans-serif; 
            max-width: 600px; 
            margin: 50px auto; 
            padding: 20px; 
            text-align: center;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 40px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .success { color: #25D366; font-size: 24px; margin-bottom: 20px; }
        .status { font-size: 18px; color: #333; margin-bottom: 30px; }
        .refresh-btn {
            background-color: #25D366;
            color: white;
            padding: 12px 24px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
            text-decoration: none;
            display: inline-block;
        }
        .refresh-btn:hover { background-color: #128C7E; }
    </style>
</head>
<body>
    <div class="container">
        <div class="success">✅ WhatsApp Bot Connected!</div>
        <div class="status">Your WhatsApp bot is connected and ready to receive messages.</div>
        <a href="/api/v1/qr/page" class="refresh-btn">Refresh Status</a>
    </div>
</body>
</html>`)
		return
	}

	_, err := h.whatsappService.GetQRCode()

	c.Header("Content-Type", "text/html")

	if err != nil {
		c.String(http.StatusOK, `
<!DOCTYPE html>
<html>
<head>
    <title>WhatsApp Bot - Setup</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { 
            font-family: Arial, sans-serif; 
            max-width: 600px; 
            margin: 50px auto; 
            padding: 20px; 
            text-align: center;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 40px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .warning { color: #f39c12; font-size: 20px; margin-bottom: 20px; }
        .message { font-size: 16px; color: #333; margin-bottom: 30px; }
        .refresh-btn {
            background-color: #3498db;
            color: white;
            padding: 12px 24px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
            text-decoration: none;
            display: inline-block;
        }
        .refresh-btn:hover { background-color: #2980b9; }
    </style>
</head>
<body>
    <div class="container">
        <div class="warning">⚠️ QR Code Not Available</div>
        <div class="message">Please wait for the WhatsApp service to generate a QR code, or refresh this page.</div>
        <div class="message">Error: %s</div>
        <a href="/api/v1/qr/page" class="refresh-btn">Refresh Page</a>
    </div>
</body>
</html>`, err.Error())
		return
	}

	c.String(http.StatusOK, `
<!DOCTYPE html>
<html>
<head>
    <title>WhatsApp Bot - QR Code Login</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { 
            font-family: Arial, sans-serif; 
            max-width: 600px; 
            margin: 50px auto; 
            padding: 20px; 
            text-align: center;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 40px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .title { color: #25D366; font-size: 24px; margin-bottom: 20px; }
        .qr-container {
            background-color: white;
            padding: 20px;
            border-radius: 10px;
            border: 2px solid #e0e0e0;
            margin: 30px 0;
            display: inline-block;
        }
        .qr-image {
            max-width: 300px;
            height: auto;
        }
        .instructions {
            text-align: left;
            background-color: #f8f9fa;
            padding: 20px;
            border-radius: 5px;
            margin: 30px 0;
        }
        .instructions h3 { color: #25D366; margin-top: 0; }
        .instructions ol { padding-left: 20px; }
        .instructions li { margin-bottom: 8px; }
        .refresh-btn {
            background-color: #25D366;
            color: white;
            padding: 12px 24px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
            text-decoration: none;
            display: inline-block;
            margin: 10px;
        }
        .refresh-btn:hover { background-color: #128C7E; }
        .auto-refresh {
            color: #666;
            font-size: 14px;
            margin-top: 20px;
        }
    </style>
    <script>
        // Auto-refresh every 30 seconds to check if connected
        setTimeout(() => {
            window.location.reload();
        }, 30000);
    </script>
</head>
<body>
    <div class="container">
        <div class="title">📱 WhatsApp Bot Setup</div>
        
        <div class="instructions">
            <h3>How to connect your WhatsApp:</h3>
            <ol>
                <li>Open WhatsApp on your phone</li>
                <li>Go to <strong>Settings</strong> → <strong>Linked Devices</strong></li>
                <li>Tap <strong>"Link a Device"</strong></li>
                <li>Scan the QR code below with your phone's camera</li>
            </ol>
        </div>
        
        <div class="qr-container">
            <img src="/api/v1/qr/image" alt="WhatsApp QR Code" class="qr-image" />
        </div>
        
        <a href="/api/v1/qr/page" class="refresh-btn">Refresh QR Code</a>
        
        <div class="auto-refresh">This page will auto-refresh every 30 seconds</div>
    </div>
</body>
</html>`)
}

func (h *qrHandler) GetQRImage(c *gin.Context) {
	qrCodeData, err := h.whatsappService.GetQRCode()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "QR code not available",
			"message": err.Error(),
		})
		return
	}

	qrPNG, err := qrcode.Encode(qrCodeData, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate QR image",
			"message": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	c.Data(http.StatusOK, "image/png", qrPNG)
}
