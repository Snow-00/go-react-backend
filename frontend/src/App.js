import React, { Fragment, useCallback, useEffect, useState } from "react";
import { Link, Outlet, useNavigate } from "react-router-dom";
import Alert from "./components/Alert";

function App() {
  const [jwtToken, setJwtToken] = useState("")
  const [alertMessage, setAlertMessage] = useState("")
  const [alertClassName, setAlertClassName] = useState("d-none")  // display property of bootstrap

  const navigate = useNavigate()

  const logout = () => {
    const requestOptions = {
      method: "GET",
      credentials: "include",
    }

    // here we dont need .then becos we just want to delete the cookie
    fetch(`${process.env.REACT_APP_BACKEND}/logout`, requestOptions)
      .catch(error => console.log("error logging out", error))
      .finally(() => {
        setJwtToken("")
      }) 
    navigate("/login")
  }

  const toggleRefresh = useCallback(() => {
    return new Promise((resolve, reject) => {
      const requestOptions = {
        method: "GET",
        credentials: "include",
      }

      fetch(`${process.env.REACT_APP_BACKEND}/refresh`, requestOptions)
        .then(response => response.json())
        .then(data => {
          if (data.error) {
            reject(data.message)
            return
          }
          resolve(data.access_token)
        })
        .catch(error => {
          reject(error)
        })
      })
  }, [])

  useEffect(() => {
    toggleRefresh()
      .then(token => {
        setJwtToken(token)
        console.log("from use effect")
      })
      .catch(error => {
        if (error.message !== "Unexpected end of JSON input") {
          console.log(error)
        }
      })
  }, [toggleRefresh])

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
            jwtToken,  
            setJwtToken,
            setAlertClassName,
            setAlertMessage,
            toggleRefresh,
          }}/>
        </div>
      </div>
    </div>
  )
}

export default App;
