import { Observable } from 'rxjs';

/**
 * data model
 */
import { EntityBean, ContentListResponse } from '../models/commons/entity-bean';

export interface DefaultResource<T extends EntityBean> {
    GetAll(): Observable<T[]>;
    GetAllFromContent(filter: string, params: {[key: string]: any | any[]}): Observable<ContentListResponse<T>>;
    GetSingle(id: string): Observable<any>;
    GetSingleAny(id: string): Observable<any>;
    Task(path: String, payload: any): Observable<any>;
    Add(itemToAdd: T): Observable<any>;
    Update(id: string, itemToUpdate: T): Observable<any>;
    Delete(id: string): Observable<any>;
}

export interface DefaultStreamResource<T> {
    GetSingle(id: string): Observable<any>;
}

export interface DefaultGraphResource<T> {
    Get(params: any): Observable<any>;
}
