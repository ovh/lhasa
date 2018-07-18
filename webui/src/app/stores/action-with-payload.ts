import { Action } from '@ngrx/store';
import { Subject } from 'rxjs';

export class ActionWithPayload<T> implements Action {
  readonly type: string;
  constructor(public payload: T) {
  }
}

export class ActionWithPayloadAndPromise<T> implements Action {
  readonly type: string;
  constructor(public payload: T, public subject?: Subject<any>) {
  }
}
