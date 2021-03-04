import React from 'react'

import { useState } from "react";
import { useHistory } from 'react-router-dom';

import firebase from "firebase/app";
import "firebase/auth";

export default function SignIn(props) {
    console.log('sign in page')
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const history = useHistory()
  
    const onSubmit = () => firebase.auth().signInWithEmailAndPassword(email, password)
    .then(userCredential => {
      firebase.auth().currentUser.getIdToken(/* forceRefresh */ true).then(idToken => {
        console.log(idToken)
      }).catch(error => {
        console.log(error)
      });
      history.push('/cms')
    })
    .catch((error) => {
      console.log(error)
      history.push('/cms')
    });
  
    return (
    <>
      <label>
        Email:
        <input type="text" value={email} onChange={e => setEmail(e.target.value)} />
      </label>
      <label>
        Password:
        <input type="text" value={password} onChange={e => setPassword(e.target.value)} />
      </label>
      <input type="submit" value="Submit" onClick={onSubmit}/>
    </>
    )
  }