import {Component} from 'angular2/core';
import {Card} from '../models/card';

@Component({
  selector: 'dos-card',
  templateUrl: 'templates/components/card.html',
})
export class CardComponent {
  @Input() card: Card;
}

