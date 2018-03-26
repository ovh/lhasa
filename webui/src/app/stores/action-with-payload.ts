import { Action } from '@ngrx/store';

export class ActionWithPayload<T> implements Action {
  readonly type: string;
  constructor(public payload: T) {
    
  }
}
