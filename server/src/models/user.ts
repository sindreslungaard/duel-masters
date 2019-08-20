import mongoose, { Schema, Document } from "mongoose"

export interface IUser extends Document {
  uid: string,
  username: string,
  email: string,
  password: string,
  sessions: [{
    token: string,
    ip: string,
    expires: number
  }],
  permissions: [string]
}
 
const User = new Schema({
  uid: String,
  username: String,
  email: String,
  password: String,
  sessions: [{
    token: String,
    ip: String,
    expires: Number
  }],
  permissions: [String]
})

export default mongoose.model<IUser>('user', User)