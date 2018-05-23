import { Observable } from 'rxjs/Observable';

/**
 * data model
 */
import { EntityBean, ContentListResponse } from '../models/commons/entity-bean';

export interface DefaultLinkResource<T extends EntityBean> {
    GetAll(id: string): Observable<T[]>;
    FindAll(id: string, filters: string): Observable<T[]>;
    GetSingle(id: string, linkId: string): Observable<T>;
    Add(id: string, linkId: string, linkToAdd: any): Observable<T>;
    Update(id: string, linkId: string, instance: string, linkToUpdate: any): Observable<T>;
    Delete(id: string, linkId: string, instance: string): Observable<T>;
    DeleteWithFilter(id: string, linkId: string, instance: string, filters: string): Observable<T>;
}

export interface DefaultResource<T extends EntityBean> {
    GetAll(): Observable<T[]>;
    GetAllFromContent(filter: string, params: {[key: string]: any | any[]}): Observable<ContentListResponse<T>>;
    GetSingle(id: string): Observable<T>;
    Add(itemToAdd: T): Observable<T>;
    Update(id: string, itemToUpdate: T): Observable<T>;
    Delete(id: string): Observable<T>;
}

export interface DefaultStreamResource<T> {
    GetSingle(id: string): Observable<T>;
}
