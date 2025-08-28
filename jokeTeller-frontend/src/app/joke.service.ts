import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Joke {
  Setup: string;
  Deliver: string;
  Type: string;
  Category: string;
  Joke: string;
  Flags: string[];
}

@Injectable({
  providedIn: 'root',
})
export class JokeService {
  private http = inject(HttpClient);
  private apiUrl = 'http://localhost:8080/api/joke';

  getJoke(category: string): Observable<Joke> {
    const url = `${this.apiUrl}?category=${category}`;
    return this.http.get<Joke>(url);
  }
}
