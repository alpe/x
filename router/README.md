# HTTP router


Router provides regular expression based path matching router. Routes are
eximined in order they are defined until match is found.

Route structure contains of three parts:

* method string being one or more, coma separated list of methods
* test string that is compiled to regular expression
* handler function

Test string can use `{name}` to define named match that will catch everything
except `/` character. It is possible to define custom regular expression rule
after `:`, for example `{id:\d+}` to match only numbers.

Example application using router:


    func main() {
        ctx := context.Background()

        // prepare context

        app := CreateApplication(ctx)
        http.ListenAndServe("localhost:8000", app)
    }

    func CreateApplication(ctx context.Context) http.Handler {
        return &application{
            ctx: ctx,
            rt: router.New(router.Routes{
                {"GET", `/users`, HandleListUsers},
                {"GET", `/users/{id}`, HandleUserDetails},
                {"POST,PUT", `/users`, HandleCreateUser},
                {"PATCH", `/user/{id:\d+}`, HandleUpdateUser},
                {router.AnyMethod, `.*`, handle404},
            }),
        }

    }

    // application is binding router with context
    type application struct {
        rt  *router.Router
        ctx context.Context
    }

    func (app *application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        app.rt.ServeCtxHTTP(app.ctx, w, r)
    }

    func handle404(ctx context.Context, w http.ResponseWriter, r *http.Request) {
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
    }

    func HandleUserDetails(ctx context.Context, w http.ResponseWriter, r *http.Request) {
        userID := web.Args(ctx).ByName("id")
        userID = web.Args(ctx).ByIndex(0)
        fmt.Fprintf(w, "details for %s", userID)
    }
