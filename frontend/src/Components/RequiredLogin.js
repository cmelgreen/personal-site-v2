import React from 'react'

import firebase from "firebase/app";
import "firebase/auth";
import {
  FirebaseAuthProvider,
  IfFirebaseAuthed,
  IfFirebaseUnAuthed
} from "@react-firebase/auth";

export default function RequiredLogn(props) {
    return (
        <FirebaseAuthProvider {...props.config} firebase={firebase}>
            <IfFirebaseAuthed>
            {({ user }) => (
                props.children[0]
            )}
            </IfFirebaseAuthed>
            <IfFirebaseUnAuthed>
                {props.children[1]}
            </IfFirebaseUnAuthed>
        </FirebaseAuthProvider>
    )
}

