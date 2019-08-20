import mongoose, { Schema, Document } from "mongoose"

export interface IDeck extends Document {
  uid: string,
  owner: string,
  name: string,
  cards: [string],
  public: boolean,
  standard: boolean
}
 
const Deck = new Schema({
  uid: String,
  owner: String,
  name: String,
  cards: [String],
  public: Boolean,
  standard: Boolean
})

export default mongoose.model<IDeck>('deck', Deck)