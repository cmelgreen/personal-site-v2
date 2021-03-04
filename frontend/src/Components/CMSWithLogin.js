import React from 'react'

import CMS from './CMS/CMS'
import SignIn from './SignIn'

import firebase from "firebase/app";
import "firebase/auth";
import {
  FirebaseAuthProvider,
  IfFirebaseAuthed,
  IfFirebaseUnAuthed
} from "@react-firebase/auth";

export default function CMSWithLogin(props) {
    return (
        <FirebaseAuthProvider {...props.config} firebase={firebase}>
            <IfFirebaseAuthed>
            {({ user }) => (
                <CMS user={user}/>
            )}
            </IfFirebaseAuthed>
            <IfFirebaseUnAuthed>
               <SignIn/>
            </IfFirebaseUnAuthed>
        </FirebaseAuthProvider>
    )
}