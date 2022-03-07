package settings

type Java struct {
	Lombok Lombok
}

type Lombok struct {
	Data               bool
	Getter             bool
	Setter             bool
	Slf4j              bool
	NoArgsConstructor  bool
	AllArgsConstructor bool
	ToString           bool
	EqualsAndHashCode  bool
}
