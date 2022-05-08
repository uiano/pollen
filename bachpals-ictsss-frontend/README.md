## Installation
Make sure to run `yarn` to install all required dependencies.

## Available Scripts

In the project directory, you can run:

### `yarn start`

Runs the app in the development mode.\
Open [http://localhost:80](http://localhost:80) to view it in the browser.

The page will reload if you make edits.\
You will also see any lint errors in the console.

### `yarn build`
Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

### Environment variables
Create file called `.env` and fill in required environment variables. Use `example.env` as an example.

#### Development environment variables
In development make sure to add following in your `.env` file.

`DANGEROUSLY_DISABLE_HOST_CHECK=true`

Also remember to replace your `API_PATH` and point it to the backend running locally.
`REACT_APP_API_PATH=http://ikt-stack.internal.uia.no/api/v1`
`REACT_APP_AUTH_PATH=http://ikt-stack.internal.uia.no/oauth2`


#### Other development values
Proxy to make react think it's performing requests on the same domain.
At the bottom in `package.json` add following.
`"proxy": "http://127.0.0.2:80/"`

In your `hosts` file, make sure your localhost points to the oauth2 domain,`frontend.ikt-stack.internal.uia.no`
The localhost ip you should point this domain to is `127.0.0.1`

From now and onwards, you must use `ikt-stack.internal.uia.no` url to reach this project, otherwise the oauth2 client will not allow you to log in.