import { Observable } from 'rxjs/Observable';

/**
 * data model
 */
import { EntityBean, ContentListResponse } from '../models/commons/entity-bean';

export interface DefaultResource<T extends EntityBean> {
    GetAll(): Observable<T[]>;
    GetAllFromContent(filter: string, params: {[key: string]: any | any[]}): Observable<ContentListResponse<T>>;
    GetSingle(id: string): Observable<T>;
    GetSingleAny(id: string): Observable<any>;
    Task(path: String, payload: any): Observable<any>;
    Add(itemToAdd: T): Observable<T>;
    Update(id: string, itemToUpdate: T): Observable<T>;
    Delete(id: string): Observable<T>;
}

export interface DefaultStreamResource<T> {
    GetSingle(id: string): Observable<T>;
}

export interface DefaultGraphResource<T> {
    Get(params: any): Observable<T>;
}
