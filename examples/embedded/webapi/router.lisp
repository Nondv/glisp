(define router
        (lambda ()
          (let ((path (alist/get "path" request-data)))
            (print (+ (alist/get "method" request-data) " " path))

            (cond
             ((= path "/hello")
              (let ((name-param (->> request-data
                                     (alist/get "query")
                                     (alist/get "name"))))
                (if name-param
                    (response 200 (+ "Hello, " name-param))
                    (response 400 "Provide `name=` parameter"))))
             ("else"
              (response 200 "It works! Try /hello"))))))


(define response
        (lambda (status body)
          (->> ()
               (alist/set "body" body)
               (alist/set "status" status))))
