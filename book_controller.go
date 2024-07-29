package controller

import (
	"context"
	"database/sql"

	libraryv1 "github.com/example/library-operator/api/v1"
	_ "github.com/lib/pq"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type BookReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	DB     *sql.DB
}

func (r *BookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	var book libraryv1.Book
	if err := r.Get(ctx, req.NamespacedName, &book); err != nil {
		log.Error(err, "unable to fetch Book")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	_, err := r.DB.Exec("INSERT INTO books (name, title, author) VALUES ($1, $2, $3)", book.Name, book.Spec.Title, book.Spec.Author)
	if err != nil {
		log.Error(err, "failed to insert book data into database")
		return ctrl.Result{}, err
	}

	log.Info("Reconciling Book", "name", book.Name)
	log.Info("Book Title", "title", book.Spec.Title)
	log.Info("Book Author", "author", book.Spec.Author)

	return ctrl.Result{}, nil
}

func (r *BookReconciler) SetupWithManager(mgr manager.Manager, db *sql.DB) error {
	r.DB = db
	return ctrl.NewControllerManagedBy(mgr).
		For(&libraryv1.Book{}).
		Complete(r)
}
