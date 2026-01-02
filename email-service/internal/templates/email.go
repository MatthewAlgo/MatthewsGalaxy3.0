package templates

import (
	"bytes"
	"html/template"
)

// NewPostEmail generates the HTML email for a new post notification
func NewPostEmail(postTitle, postExcerpt, postURL, unsubscribeURL string) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Post from Matthew's Galaxy</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            background: linear-gradient(135deg, #0a0a1a 0%, #1a0a2e 100%);
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            color: #ffffff;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 40px 20px;
        }
        .header {
            text-align: center;
            padding: 30px 0;
        }
        .stars {
            font-size: 24px;
            letter-spacing: 10px;
        }
        .logo {
            font-size: 28px;
            font-weight: bold;
            background: linear-gradient(90deg, #4a9eff, #ff6b9d);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }
        .content {
            background: rgba(255, 255, 255, 0.05);
            border-radius: 16px;
            padding: 30px;
            margin: 20px 0;
            border: 1px solid rgba(255, 255, 255, 0.1);
        }
        .post-title {
            font-size: 24px;
            color: #4a9eff;
            margin-bottom: 16px;
            text-decoration: none;
        }
        .post-title a {
            color: #4a9eff;
            text-decoration: none;
        }
        .post-excerpt {
            font-size: 16px;
            line-height: 1.6;
            color: #b8b8d0;
            margin-bottom: 24px;
        }
        .cta-button {
            display: inline-block;
            background: linear-gradient(90deg, #4a9eff, #ff6b9d);
            color: #ffffff;
            padding: 14px 32px;
            border-radius: 50px;
            text-decoration: none;
            font-weight: bold;
            font-size: 16px;
        }
        .footer {
            text-align: center;
            padding: 30px 0;
            font-size: 14px;
            color: #666;
        }
        .unsubscribe {
            color: #888;
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="stars">✨ 🌟 ⭐ 🌟 ✨</div>
            <h1 class="logo">Matthew's Galaxy</h1>
        </div>
        
        <div class="content">
            <h2 style="color: #ffffff; margin-bottom: 8px;">New Post Alert! 🚀</h2>
            <p style="color: #b8b8d0; margin-bottom: 24px;">Matei just published something new:</p>
            
            <h3 class="post-title"><a href="{{.PostURL}}">{{.PostTitle}}</a></h3>
            <p class="post-excerpt">{{.PostExcerpt}}</p>
            
            <a href="{{.PostURL}}" class="cta-button">Read the Full Post →</a>
        </div>
        
        <div class="footer">
            <p>You're receiving this because you subscribed to Matthew's Galaxy.</p>
            <p><a href="{{.UnsubscribeURL}}" class="unsubscribe">Unsubscribe</a></p>
            <p>© 2026 Matthew's Galaxy. Ad astra per aspera ✨</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("newpost").Parse(tmpl)
	if err != nil {
		return "", err
	}

	data := struct {
		PostTitle      string
		PostExcerpt    string
		PostURL        string
		UnsubscribeURL string
	}{
		PostTitle:      postTitle,
		PostExcerpt:    postExcerpt,
		PostURL:        postURL,
		UnsubscribeURL: unsubscribeURL,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WelcomeEmail generates the HTML email for new subscribers
func WelcomeEmail(userName, unsubscribeURL string) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Matthew's Galaxy</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            background: linear-gradient(135deg, #0a0a1a 0%, #1a0a2e 100%);
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            color: #ffffff;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 40px 20px;
        }
        .header {
            text-align: center;
            padding: 30px 0;
        }
        .stars {
            font-size: 32px;
            letter-spacing: 10px;
        }
        .logo {
            font-size: 32px;
            font-weight: bold;
            background: linear-gradient(90deg, #4a9eff, #ff6b9d);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }
        .content {
            background: rgba(255, 255, 255, 0.05);
            border-radius: 16px;
            padding: 40px;
            margin: 20px 0;
            border: 1px solid rgba(255, 255, 255, 0.1);
            text-align: center;
        }
        .welcome-text {
            font-size: 28px;
            color: #ffffff;
            margin-bottom: 16px;
        }
        .description {
            font-size: 16px;
            line-height: 1.8;
            color: #b8b8d0;
            margin-bottom: 24px;
        }
        .footer {
            text-align: center;
            padding: 30px 0;
            font-size: 14px;
            color: #666;
        }
        .unsubscribe {
            color: #888;
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="stars">🌌 ✨ 🚀 ✨ 🌌</div>
            <h1 class="logo">Matthew's Galaxy</h1>
        </div>
        
        <div class="content">
            <h2 class="welcome-text">Welcome aboard, {{.UserName}}! 🎉</h2>
            <p class="description">
                You've successfully subscribed to Matthew's Galaxy newsletter!
                <br><br>
                You'll receive updates whenever Matei publishes new content about:
                <br>
                ⚡ Cloud Architecture & Distributed Systems
                <br>
                💻 Software Engineering Insights
                <br>
                🚀 Project Showcases & Tutorials
                <br>
                ✨ Career & Tech Reflections
            </p>
            <p style="font-size: 24px;">Ad astra per aspera! 🌟</p>
        </div>
        
        <div class="footer">
            <p>You're receiving this because you subscribed to Matthew's Galaxy.</p>
            <p><a href="{{.UnsubscribeURL}}" class="unsubscribe">Unsubscribe</a></p>
            <p>© 2026 Matthew's Galaxy</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("welcome").Parse(tmpl)
	if err != nil {
		return "", err
	}

	data := struct {
		UserName       string
		UnsubscribeURL string
	}{
		UserName:       userName,
		UnsubscribeURL: unsubscribeURL,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
