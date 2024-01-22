import React, { Fragment, useState } from "react";
import { Link, Outlet, useNavigate } from "react-router-dom";
import Alert from "./components/Alert";

function App() {
  const [jwtToken, setJwtToken] = useState("")
  const [alertMessage, setAlertMessage] = useState("")
  const [alertClassName, setAlertClassName] = useState("d-none")  // display property of bootstrap

  const navigate = useNavigate()

  const logout = () => {
    setJwtToken("")
    navigate("/login")
  }

  return (
    <div className="container">
      <div className="row">
        <div className="col">
          <h1 className="mt-3">Go Watch a Movie</h1>
        </div>

        <div className="col text-end">
          {jwtToken === ""
            ? <Link to="/login"><span className="badge bg-success">Login</span></Link>
            : <a href="#" onClick={logout}><span className="badge bg-danger">Logout</span></a>
          }
        </div>
        <hr className="mb-3"/>
      </div>

      <div className="row">
        <div className="col-md-2">
          <nav>
            <div className="list-group">
              <Link to="/" className="list-group-item list-group-action">Home</Link>
              <Link to="/movies" className="list-group-item list-group-action">Movies</Link>
              <Link to="/genres" className="list-group-item list-group-action">Genres</Link>
              {jwtToken !== "" &&
                <Fragment>
                  <Link to="/admin/movie/0" className="list-group-item list-group-action">Add Movies</Link>
                  <Link to="/manage-catalog" className="list-group-item list-group-action">Manage Catalog</Link>
                  <Link to="/graphql" className="list-group-item list-group-action">GraphQL</Link>
                </Fragment>
              }
            </div>
          </nav>
        </div>

        <div className="col-md-10">
          <Alert message={alertMessage} className={alertClassName} />
          <Outlet context={{
            // jwtToken,  -> if doesnt need to pass the jwtToken, this isnt needed here
            setJwtToken,
            setAlertClassName,
            setAlertMessage,
          }}/>
        </div>
      </div>
    </div>
  )
}

export default App;
