import React, { Fragment, useCallback, useEffect, useState } from "react";
import { Link, Outlet, useNavigate } from "react-router-dom";
import Alert from "./components/Alert";

function App() {
  const [jwtToken, setJwtToken] = useState("")
  const [alertMessage, setAlertMessage] = useState("")
  const [alertClassName, setAlertClassName] = useState("d-none")  // display property of bootstrap

  const [tickInterval, setTickInterval] = useState()

  const navigate = useNavigate()

  const logout = () => {
    const requestOptions = {
      method: "GET",
      credentials: "include",
    }

    // here we dont need .then becos we just want to delete the cookie
    // fetch(`https://supreme-halibut-v664446pgxqxhwxvr-8080.app.github.dev/logout`, requestOptions)
    fetch(`http://localhost:8080/logout`, requestOptions)
      .catch(error => console.log("error logging out", error))
      .finally(() => {
        setJwtToken("")
        // toggleRefresh(false)
      }) 
    navigate("/login")
  }

  // const toggleRefresh = useCallback(status => {
  //   console.log("clicked")

  //   if (status) {
  //     let i = setInterval(() => {
  //       console.log("this will run every sec")
  //     }, 1000)
  //     setTickInterval(i)

  //     console.log("setting tick interval to", i)
  //   } else {
  //     console.log("turning off tickInterval", tickInterval)
      
  //     setTickInterval(null)
  //     clearInterval(tickInterval)
  //   }
  // }, [tickInterval])

  const toggleRefresh = useCallback(() => {
    if (jwtToken !== "") {
      setInterval
    }
  })

  useEffect(() => {
    console.log("use effect run")

    if (jwtToken === "") {
      console.log("the jwt token:", jwtToken)
      
      const requestOptions = {
        method: "GET",
        credentials: "include",
      }

      // fetch(`https://supreme-halibut-v664446pgxqxhwxvr-8080.app.github.dev/refresh`, requestOptions)
      fetch(`http://localhost:8080/refresh`, requestOptions)
        .then(response => response.json())
        .then(data => {
          if (data.access_token) {
            setJwtToken(data.access_token)
            // toggleRefresh(true)
          }
        })
        .catch(error => console.log("user not logged in", error))
    }
  // }, [jwtToken, refreshToggle]
  }, [jwtToken])

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
            toggleRefresh,
          }}/>
        </div>
      </div>
    </div>
  )
}

export default App;
