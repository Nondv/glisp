(load "lang/core.lisp")

(define alist/get
        (lambda (key alist)
          (if (not alist)
              ()
              (let ((cell (car alist)))
                (if (= key (car cell))
                    (cdr cell)
                    (alist/get key (cdr alist)))))))

(define alist/set
        (lambda (key value alist)
          (cons (cons key value)
                alist)))
