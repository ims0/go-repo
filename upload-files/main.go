package main

import (
    "crypto/md5"
    "fmt"
    "io"
    "mime/multipart"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/kataras/iris/v12"
)

// SaveFormFile saves a result of `FormFile` to the "dest" disk full path (directory + filename).
// See `FormFile` and `UploadFormFiles` too.
func  SaveFormFile(fh *multipart.FileHeader, dest string) (int64, error) {
    src, err := fh.Open()
    if err != nil {
        return 0, err
    }
    defer src.Close()

    out, err := os.Create(dest)
    if err != nil {
        return 0, err
    }
    defer out.Close()

    return io.Copy(out, src)
}

func main() {
    app := newApp()
    // start the server at http://localhost:8080 with post limit at 32 MB.
    app.Listen(":8080",  iris.WithRemoteAddrHeader("X-Forwarded-For"))
}

func newApp() *iris.Application {
    app := iris.New()
    app.RegisterView(iris.HTML("./view", ".html"))

    // Serve the index.html to the client.
    app.Get("/file", func(ctx iris.Context) {
        // create a token (optionally).

        now := time.Now().Unix()
        h := md5.New()
        io.WriteString(h, strconv.FormatInt(now, 10))
        token := fmt.Sprintf("%x", h.Sum(nil))

        // render the form with the token for any use you'd like.
        ctx.View("index.html", token)
    })

    // Handle the post request from the index.html to the server.
    app.Post("/file", func(ctx iris.Context) {
        //
        // UploadFormFiles
        // uploads any number of incoming files ("multiple" property on the form input).
        //

        // second argument is optional,
        // it can be used to change a file's name based on the request,
        // at this example we will showcase how to use it
        // by prefixing the uploaded file with the current user's ip.
        ctx.UploadFormFiles("./share", beforeSave)
        ctx.View("index.html", 1)
    })

    app.Post("/upload_manual", func(ctx iris.Context) {
        // Get the max post value size passed via iris.WithPostMaxMemory.
        maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()

        err := ctx.Request().ParseMultipartForm(maxSize)
        if err != nil {
            ctx.StopExecution()
            return
        }

        form := ctx.Request().MultipartForm

        files := form.File["files[]"]
        failures := 0
        for _, file := range files {
            _, err = SaveFormFile(file, "./uploads/"+file.Filename)
            if err != nil {
                failures++
                ctx.Writef("failed to upload: %s\n", file.Filename)
            }
        }
        ctx.Writef("%d files uploaded", len(files)-failures)
    })

    return app
}

func beforeSave(ctx iris.Context, file *multipart.FileHeader)  {
    ip := ctx.RemoteAddr()
    fmt.Println("upload from host:",ip)
    // make sure you format the ip in a way
    // that can be used for a file name (simple case):
    ip = strings.ReplaceAll(ip, ".", "_")
    ip = strings.ReplaceAll(ip, ":", "_")

    // you can use the time.Now, to prefix or suffix the files
    // based on the current time as well, as an exercise.
    // i.e unixTime :=    time.Now().Unix()
    // prefix the Filename with the $IP-
    // no need for more actions, internal uploader will use this
    // name to save the file into the "./uploads" folder.
    if ip == "" {
        return  // don't change the file but continue saving it.
    }

    file.Filename = ip + "-" + file.Filename
    return
}
