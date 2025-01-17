/*--------------------------------------------------------------------------
 * Copyright (c) 2019-2021, Postgres.ai, Nikolay Samokhvalov nik@postgres.ai
 * All Rights Reserved. Proprietary and confidential.
 * Unauthorized copying of this file, via any medium is strictly prohibited
 *--------------------------------------------------------------------------
 */

import { Component } from 'react'
import { BrowserRouter as Router, Route } from 'react-router-dom'
import {
  createGenerateClassName,
  StylesProvider,
  ThemeProvider,
} from '@material-ui/core/styles'

import { ROUTES } from 'config/routes'

import { IndexPageWrapper } from 'components/IndexPage/IndexPageWrapper'
import { theme } from '@postgres.ai/shared/styles/theme'

class App extends Component {
  render() {
    const generateClassName = createGenerateClassName({
      productionPrefix: 'p',
    })
    return (
      <Router basename={process.env.PUBLIC_URL}>
        <Route path={ROUTES.ROOT.path}>
          <StylesProvider generateClassName={generateClassName}>
            <ThemeProvider theme={theme}>
              <IndexPageWrapper />
            </ThemeProvider>
          </StylesProvider>
        </Route>
      </Router>
    )
  }
}

export default App
