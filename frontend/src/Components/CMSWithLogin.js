import React from 'react'

import RequiredLogin from './RequiredLogin'
import CMS from './CMS/CMS'
import SignIn from './SignIn'

export default function CMSWithLogin(props) {
    return (
        <RequiredLogin config={props.config}>
            <CMS/>
            <SignIn/>
        </RequiredLogin>
    )
}