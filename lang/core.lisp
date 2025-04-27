;;
;; Just some useful functions that could be a part of the language by default
;;

(define quote (lambda quote-ARG (car quote-ARG)))

(define not (lambda (x) (if x () (quote true))))

(define reduce
        (lambda (initial-value f lst)
          (if (not lst)
              initial-value
              (reduce (f initial-value (car lst)) f (cdr lst)))))

(define mapcar
        (lambda (f lst)
          (if (not lst)
              lst
              (cons (f (car lst)) (mapcar f (cdr lst))))))

(define list
        (lambda list-ARG
          (mapcar eval list-ARG)))

(define reduce
        (lambda (initial-value f lst)
          (if (not lst)
              initial-value
              (reduce (f initial-value (car lst)) f (cdr lst)))))

(define push-last
        (lambda (x lst)
          (cons (car lst)
                (if (cdr lst)
                    (push-last x (cdr lst))
                    (list x)))))

(define ->>
        (lambda ->>ARGS
          (eval (reduce (car ->>ARGS)
                        push-last
                        (cdr ->>ARGS)))))
